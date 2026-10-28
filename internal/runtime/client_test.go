package runtime

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestRuntimeHelper(t *testing.T) {
	if os.Getenv("FACTORYCTL_RUNTIME_HELPER") != "1" {
		return
	}
	mode := os.Getenv("FACTORYCTL_RUNTIME_MODE")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 4096), maxMessage)
	write := func(v any) { b, _ := json.Marshal(v); fmt.Println(string(b)) }
	for scanner.Scan() {
		var q packet
		if json.Unmarshal(scanner.Bytes(), &q) != nil {
			os.Exit(2)
		}
		if q.Method == "initialized" {
			continue
		}
		if q.Method == "" {
			write(map[string]any{"method": "reply/observed", "params": map[string]any{"id": q.ID, "result": q.Result}})
			continue
		}
		var result any = map[string]any{}
		switch q.Method {
		case "initialize":
			if mode == "init-error" {
				write(map[string]any{"id": q.ID, "error": map[string]any{"code": -32601, "message": "unsupported initialize"}})
				continue
			}
			result = map[string]any{"userAgent": "fake-runtime/1", "platformOs": "test", "platformFamily": "test"}
		case "model/list":
			result = map[string]any{"data": []any{map[string]any{"id": "qualified", "model": "qualified", "isDefault": true, "supportedReasoningEfforts": []any{map[string]any{"reasoningEffort": "low"}}}}}
		case "echo":
			write(map[string]any{"method": "item/agentMessage/delta", "params": map[string]any{"delta": "stream"}})
			result = json.RawMessage(q.Params)
		case "ask":
			write(map[string]any{"id": "native-question", "method": "item/tool/requestUserInput", "params": map[string]any{"questions": []any{}}})
		case "ignore":
			continue
		case "malformed":
			fmt.Println("not JSON")
			continue
		case "flood":
			for i := 0; i < eventCapacity+100; i++ {
				write(map[string]any{"method": "event", "params": map[string]int{"i": i}})
			}
		case "rpc-error":
			write(map[string]any{"id": q.ID, "error": map[string]any{"code": -32602, "message": "bad input"}})
			continue
		}
		write(map[string]any{"id": q.ID, "result": result})
	}
	os.Exit(0)
}
func helper(t *testing.T, cfg Config, mode string) (*Client, error) {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestRuntimeHelper$")
	cmd.Env = append(os.Environ(), "FACTORYCTL_RUNTIME_HELPER=1", "FACTORYCTL_RUNTIME_MODE="+mode)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	c, err := start(ctx, cfg, cmd)
	if c != nil {
		t.Cleanup(func() {
			if err := c.Close(); err != nil {
				t.Error(err)
			}
		})
	}
	return c, err
}
func TestInvocationOverridesAreArgumentsNotShell(t *testing.T) {
	cfg := Config{MCPCommand: "/tmp/program name", MCPArgs: []string{"mcp", "--path", "/tmp/quote\" and $(touch nope)"}, MCPEnv: map[string]string{"TOKEN": "literal\n\"secret", "URL": "http://127.0.0.1:1"}}
	args, err := Arguments(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(args) != 7 || args[0] != "app-server" || args[1] != "-c" {
		t.Fatal(args)
	}
	if !strings.Contains(args[4], "$(touch nope)") || !strings.Contains(args[6], `"TOKEN"="literal\n\"secret"`) {
		t.Fatal("values were not preserved as quoted literals", args)
	}
	if _, err = Arguments(Config{MCPCommand: "relative"}); err == nil {
		t.Fatal("relative executable accepted")
	}
	if _, err = Arguments(Config{MCPArgs: []string{"x"}}); err == nil {
		t.Fatal("orphan args accepted")
	}
}
func TestConcurrentRequestsAndNativeReplyCorrelation(t *testing.T) {
	c, err := helper(t, Config{}, "")
	if err != nil {
		t.Fatal(err)
	}
	if c.Info().UserAgent != "fake-runtime/1" {
		t.Fatal(c.Info())
	}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, e := c.Request(ctx, "echo", map[string]int{"i": i})
			if e != nil {
				t.Error(e)
				return
			}
			var got map[string]int
			json.Unmarshal(r, &got)
			if got["i"] != i {
				t.Errorf("misrouted response %s", r)
			}
		}(i)
	}
	wg.Wait()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if _, err = c.Request(ctx, "ask", nil); err != nil {
		t.Fatal(err)
	}
	var event Event
	for {
		select {
		case e := <-c.Events():
			if e.Method == "item/tool/requestUserInput" {
				event = e
			}
		case <-ctx.Done():
			t.Fatal("native request missing")
		}
		if event.Method != "" {
			break
		}
	}
	if err = c.Reply(event.ID, map[string]any{"answers": map[string]any{}}); err != nil {
		t.Fatal(err)
	}
	if err = c.Reply(event.ID, nil); err == nil {
		t.Fatal("duplicate reply accepted")
	}
	if err = c.Reply(json.RawMessage(`"other"`), nil); err == nil {
		t.Fatal("unknown reply accepted")
	}
}
func TestRequestTimeoutDoesNotKillOtherRequests(t *testing.T) {
	c, err := helper(t, Config{}, "")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	if _, err = c.Request(ctx, "ignore", nil); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("got %v", err)
	}
	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second)
	defer cancel2()
	if _, err = c.Request(ctx2, "echo", map[string]string{"value": "still alive"}); err != nil {
		t.Fatal(err)
	}
	if _, err = c.Request(ctx2, "rpc-error", nil); err == nil {
		t.Fatal("RPC error lost")
	} else {
		var rpc *RPCError
		if !errors.As(err, &rpc) || rpc.Code != -32602 {
			t.Fatal(err)
		}
	}
}
func TestRuntimeProtocolFailuresAreExplicit(t *testing.T) {
	for _, method := range []string{"malformed", "flood"} {
		t.Run(method, func(t *testing.T) {
			c, err := helper(t, Config{}, "")
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, _ = c.Request(ctx, method, nil)
			select {
			case <-c.Done():
				if c.Err() == nil {
					t.Fatal("missing failure")
				}
			case <-ctx.Done():
				t.Fatal("protocol failure did not disconnect")
			}
		})
	}
	if c, err := helper(t, Config{}, "init-error"); err == nil || c != nil {
		t.Fatal("unsupported startup accepted")
	}
}
func TestModelAndEffortRequireAdvertisedSupport(t *testing.T) {
	for _, cfg := range []Config{{Model: "unavailable"}, {Model: "qualified", Effort: "ultra"}, {Effort: "unknown"}} {
		if c, err := helper(t, cfg, ""); err == nil || c != nil {
			t.Fatal("unsupported selection accepted", cfg)
		}
	}
	if _, err := helper(t, Config{Model: "qualified", Effort: "low"}, ""); err != nil {
		t.Fatal(err)
	}
}
