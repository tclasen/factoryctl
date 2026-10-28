// Package runtime connects an explicitly selected installed Codex app-server.
// It owns no project workflow, transcript store or credentials. Personal chat settings are explicit host-owned configuration.
package runtime

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const maxMessage = 4 << 20
const eventCapacity = 1024

type Config struct {
	Executable, Dir, Model, Effort, MCPCommand string
	MCPArgs                                    []string
	MCPEnv                                     map[string]string
}
type Event struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
	ID     json.RawMessage `json:"id,omitempty"`
}
type Info struct {
	UserAgent      string `json:"user_agent"`
	PlatformOS     string `json:"platform_os"`
	PlatformFamily string `json:"platform_family"`
}
type RPCError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (e *RPCError) Error() string { return fmt.Sprintf("runtime RPC %d: %s", e.Code, e.Message) }

type packet struct {
	ID     json.RawMessage `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
	Result json.RawMessage `json:"result"`
	Error  *RPCError       `json:"error"`
}
type response struct {
	result json.RawMessage
	err    error
}
type Client struct {
	cmd            *exec.Cmd
	stdin          io.WriteCloser
	stdout         io.ReadCloser
	writeMu        sync.Mutex
	mu             sync.Mutex
	pending        map[string]chan response
	serverRequests map[string]struct{}
	sequence       uint64
	err            error
	done           chan struct{}
	processDone    chan struct{}
	readDone       chan struct{}
	events         chan Event
	stopOnce       sync.Once
	info           Info
}

// Arguments serializes only explicit invocation overrides. Values are TOML
// strings/arrays, passed as separate argv entries, never shell-interpolated.
func Arguments(cfg Config) ([]string, error) {
	args := []string{"app-server"}
	if cfg.MCPCommand == "" {
		if len(cfg.MCPArgs) > 0 || len(cfg.MCPEnv) > 0 {
			return nil, errors.New("MCP command required for MCP options")
		}
		return args, nil
	}
	if !filepath.IsAbs(cfg.MCPCommand) {
		return nil, errors.New("MCP executable must be an absolute path")
	}
	quote := func(s string) string { b, _ := json.Marshal(s); return string(b) }
	args = append(args, "-c", "mcp_servers.factoryctl.command="+quote(cfg.MCPCommand))
	av := make([]string, len(cfg.MCPArgs))
	for i, a := range cfg.MCPArgs {
		av[i] = quote(a)
	}
	args = append(args, "-c", "mcp_servers.factoryctl.args=["+strings.Join(av, ",")+"]")
	if len(cfg.MCPEnv) > 0 {
		keys := make([]string, 0, len(cfg.MCPEnv))
		for k := range cfg.MCPEnv {
			if k == "" || strings.ContainsAny(k, "=\x00") {
				return nil, errors.New("invalid MCP environment name")
			}
			keys = append(keys, k)
		}
		sort.Strings(keys)
		pairs := make([]string, 0, len(keys))
		for _, k := range keys {
			pairs = append(pairs, quote(k)+"="+quote(cfg.MCPEnv[k]))
		}
		args = append(args, "-c", "mcp_servers.factoryctl.env={"+strings.Join(pairs, ",")+"}")
	}
	return args, nil
}

// Start initializes the native protocol. The caller owns its lifetime context;
// individual Request contexts can expire without cancelling the whole client.
func Start(ctx context.Context, cfg Config) (*Client, error) {
	if cfg.Executable == "" {
		return nil, errors.New("runtime executable must be explicitly selected")
	}
	if !filepath.IsAbs(cfg.Dir) {
		return nil, errors.New("runtime directory must be absolute")
	}
	st, err := os.Stat(cfg.Dir)
	if err != nil || !st.IsDir() {
		return nil, errors.New("runtime directory is unavailable")
	}
	args, err := Arguments(cfg)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(cfg.Executable, args...)
	cmd.Dir = cfg.Dir
	return start(ctx, cfg, cmd)
}
func start(ctx context.Context, cfg Config, cmd *exec.Cmd) (*Client, error) {
	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		in.Close()
		return nil, err
	}
	// Diagnostics may contain native environment or account details. Discard them;
	// actionable protocol errors still travel through the bounded JSON-RPC path.
	cmd.Stderr = io.Discard
	c := &Client{cmd: cmd, stdin: in, stdout: out, pending: map[string]chan response{}, serverRequests: map[string]struct{}{}, done: make(chan struct{}), processDone: make(chan struct{}), readDone: make(chan struct{}), events: make(chan Event, eventCapacity)}
	if err = cmd.Start(); err != nil {
		in.Close()
		out.Close()
		return nil, fmt.Errorf("runtime unavailable: %w", err)
	}
	go c.read()
	go func() {
		<-c.readDone
		err := cmd.Wait()
		if err == nil {
			err = io.EOF
		}
		c.stop(fmt.Errorf("runtime disconnected: %w", err))
		close(c.processDone)
	}()
	go func() {
		select {
		case <-ctx.Done():
			c.stop(ctx.Err())
		case <-c.done:
		}
	}()
	initCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	raw, err := c.Request(initCtx, "initialize", map[string]any{"clientInfo": map[string]string{"name": "factoryctl", "version": "0.1.0"}, "capabilities": map[string]any{"experimentalApi": true}})
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("runtime initialization failed: %w", err)
	}
	var info struct {
		UserAgent      string `json:"userAgent"`
		PlatformOS     string `json:"platformOs"`
		PlatformFamily string `json:"platformFamily"`
	}
	if err = json.Unmarshal(raw, &info); err != nil {
		c.Close()
		return nil, errors.New("invalid runtime initialization response")
	}
	c.info = Info{UserAgent: info.UserAgent, PlatformOS: info.PlatformOS, PlatformFamily: info.PlatformFamily}
	if err = c.write(map[string]any{"method": "initialized"}); err != nil {
		c.Close()
		return nil, err
	}
	if cfg.Model != "" || cfg.Effort != "" {
		if err = c.validateSelection(initCtx, cfg); err != nil {
			c.Close()
			return nil, err
		}
	}
	return c, nil
}
func (c *Client) Info() Info            { return c.info }
func (c *Client) Events() <-chan Event  { return c.events }
func (c *Client) Done() <-chan struct{} { return c.done }
func (c *Client) Err() error            { c.mu.Lock(); defer c.mu.Unlock(); return c.err }
func (c *Client) write(v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if len(b) > maxMessage {
		return errors.New("runtime message exceeds size limit")
	}
	b = append(b, '\n')
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	select {
	case <-c.done:
		return c.Err()
	default:
	}
	timer := time.AfterFunc(10*time.Second, func() { c.stop(errors.New("runtime write timed out")) })
	defer timer.Stop()
	_, err = c.stdin.Write(b)
	if err != nil {
		c.stop(fmt.Errorf("runtime write failed: %w", err))
	}
	return err
}
func (c *Client) Request(ctx context.Context, method string, params any) (json.RawMessage, error) {
	if method == "" {
		return nil, errors.New("runtime method is empty")
	}
	c.mu.Lock()
	if c.err != nil {
		err := c.err
		c.mu.Unlock()
		return nil, err
	}
	if len(c.pending) >= eventCapacity {
		c.mu.Unlock()
		return nil, errors.New("too many pending runtime requests")
	}
	c.sequence++
	id := strconv.FormatUint(c.sequence, 10)
	reply := make(chan response, 1)
	c.pending[id] = reply
	c.mu.Unlock()
	defer func() { c.mu.Lock(); delete(c.pending, id); c.mu.Unlock() }()
	if err := c.write(map[string]any{"id": json.RawMessage(id), "method": method, "params": params}); err != nil {
		return nil, err
	}
	select {
	case r := <-reply:
		return r.result, r.err
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.done:
		select {
		case r := <-reply:
			return r.result, r.err
		default:
			return nil, c.Err()
		}
	}
}

// Reply answers only an outstanding runtime-originated request. The trusted host
// must validate the user's selected response against that request's schema.
func (c *Client) Reply(id json.RawMessage, result any) error {
	key := string(id)
	c.mu.Lock()
	_, ok := c.serverRequests[key]
	if ok {
		delete(c.serverRequests, key)
	}
	c.mu.Unlock()
	if !ok {
		return errors.New("unknown or already answered runtime request")
	}
	if err := c.write(map[string]any{"id": id, "result": result}); err != nil {
		return err
	}
	return nil
}
func (c *Client) read() {
	defer close(c.readDone)
	defer close(c.events)
	scanner := bufio.NewScanner(c.stdout)
	scanner.Buffer(make([]byte, 4096), maxMessage)
	for scanner.Scan() {
		var p packet
		if err := json.Unmarshal(scanner.Bytes(), &p); err != nil {
			c.stop(errors.New("runtime emitted malformed protocol JSON"))
			return
		}
		if p.Method == "" {
			if len(p.ID) == 0 {
				c.stop(errors.New("runtime response has no identity"))
				return
			}
			c.mu.Lock()
			ch := c.pending[string(p.ID)]
			c.mu.Unlock()
			if ch != nil {
				r := response{result: p.Result}
				if p.Error != nil {
					r.err = p.Error
				}
				select {
				case ch <- r:
				default:
				}
			}
			continue
		}
		if len(p.ID) > 0 && string(p.ID) != "null" {
			c.mu.Lock()
			if len(c.serverRequests) >= eventCapacity {
				c.mu.Unlock()
				c.stop(errors.New("too many unanswered native runtime requests"))
				return
			}
			c.serverRequests[string(p.ID)] = struct{}{}
			c.mu.Unlock()
		}
		select {
		case c.events <- Event{Method: p.Method, Params: p.Params, ID: p.ID}:
		case <-c.done:
			return
		default:
			c.stop(errors.New("runtime event buffer exceeded; conversation disconnected rather than dropping events"))
			return
		}
	}
	if err := scanner.Err(); err != nil {
		c.stop(fmt.Errorf("runtime protocol read failed: %w", err))
	} else {
		c.stop(io.EOF)
	}
}
func (c *Client) stop(err error) {
	c.stopOnce.Do(func() {
		c.mu.Lock()
		c.err = err
		c.mu.Unlock()
		close(c.done)
		c.stdin.Close()
		if c.cmd.Process != nil {
			_ = c.cmd.Process.Kill()
		}
		c.stdout.Close()
	})
}
func (c *Client) Close() error {
	c.stop(errors.New("runtime closed"))
	select {
	case <-c.processDone:
		return nil
	case <-time.After(3 * time.Second):
		return errors.New("runtime shutdown timed out")
	}
}

func (c *Client) validateSelection(ctx context.Context, cfg Config) error {
	models, err := c.Models(ctx)
	if err != nil {
		return err
	}
	_, err = ValidateChatSelection(ChatSelection{Model: cfg.Model, Effort: cfg.Effort}, models)
	return err
}
