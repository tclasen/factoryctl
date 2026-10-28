// Package mcp exposes scoped document operations and interaction requests over
// newline-delimited JSON-RPC stdio. It never reads stdin as human answers.
package mcp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/tclasen/factoryctl/internal/app"
	"github.com/tclasen/factoryctl/internal/knowledge"
	"github.com/tclasen/factoryctl/internal/project"
)

type Config struct {
	Projects    []project.Project
	Selected    project.Project
	BridgeURL   string
	BridgeToken string
}
type Server struct {
	config Config
	client *http.Client
}
type Request struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}
type BridgeResponse struct {
	Result json.RawMessage `json:"result,omitempty"`
	Error  *app.Error      `json:"error,omitempty"`
}

func New(config Config) (*Server, error) {
	if config.BridgeURL != "" {
		u, e := url.Parse(config.BridgeURL)
		if e != nil || u.Scheme != "http" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || !loopback(u.Hostname()) || config.BridgeToken == "" {
			return nil, app.Fail("invalid_arguments", "interaction bridge must be an authenticated loopback HTTP endpoint")
		}
	}
	if config.Selected.Checkout != "" {
		found := false
		for _, p := range config.Projects {
			if p.Checkout == config.Selected.Checkout && p.Root == config.Selected.Root {
				found = true
			}
		}
		if !found {
			return nil, app.Fail("invalid_arguments", "selected project is outside MCP scope")
		}
	}
	return &Server{config: config, client: &http.Client{Timeout: 35 * time.Second, Transport: &http.Transport{Proxy: nil}, CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }}}, nil
}
func loopback(host string) bool { ip := net.ParseIP(host); return ip != nil && ip.IsLoopback() }
func strict(data json.RawMessage, v any) error {
	if len(data) == 0 {
		data = []byte("{}")
	}
	if len(bytes.TrimSpace(data)) == 0 || bytes.TrimSpace(data)[0] != '{' {
		return app.Fail("invalid_arguments", "expected an object")
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if e := d.Decode(v); e != nil {
		return app.Fail("invalid_arguments", e.Error())
	}
	if e := d.Decode(new(any)); e != io.EOF {
		return app.Fail("invalid_arguments", "expected one JSON object")
	}
	return nil
}
func (s *Server) Call(ctx context.Context, name string, args json.RawMessage) (any, error) {
	switch name {
	case "factory_user_ask", "factory_user_wait", "factory_ui_present":
		return s.bridge(ctx, name, args)
	}
	return s.documentCall(name, args)
}
func (s *Server) bridge(ctx context.Context, name string, args json.RawMessage) (any, error) {
	if s.config.BridgeURL == "" || s.config.BridgeToken == "" {
		return nil, app.Fail("interaction_unavailable", "no trusted user surface is attached; use a web or interactive conversation")
	}
	body, e := json.Marshal(Request{Name: name, Arguments: args})
	if e != nil {
		return nil, e
	}
	req, e := http.NewRequestWithContext(ctx, http.MethodPost, s.config.BridgeURL, bytes.NewReader(body))
	if e != nil {
		return nil, e
	}
	req.Header.Set("Authorization", "Bearer "+s.config.BridgeToken)
	req.Header.Set("Content-Type", "application/json")
	resp, e := s.client.Do(req)
	if e != nil {
		return nil, app.Fail("interaction_unavailable", "attached conversation is unavailable")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, app.Fail("interaction_unavailable", fmt.Sprintf("attached surface returned HTTP %d", resp.StatusCode))
	}
	data, e := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if e != nil {
		return nil, e
	}
	var result BridgeResponse
	if e = json.Unmarshal(data, &result); e != nil {
		return nil, app.Fail("interaction_unavailable", "invalid bridge response")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	if len(result.Result) == 0 {
		return nil, app.Fail("interaction_unavailable", "bridge omitted its result")
	}
	return result.Result, nil
}

type nodeArgs struct {
	Ref              knowledge.Ref  `json:"ref"`
	ExpectedRevision string         `json:"expected_revision,omitempty"`
	Set              map[string]any `json:"set,omitempty"`
	Document         *string        `json:"document,omitempty"`
}

func (s *Server) selected(ref knowledge.Ref) (project.Project, error) {
	p := s.config.Selected
	if p.Checkout == "" {
		if len(s.config.Projects) != 1 {
			return p, app.Fail("ambiguous_project", "MCP requires an explicit project selection")
		}
		p = s.config.Projects[0]
	}
	if ref.ProjectID != "" && ref.ProjectID != p.ID || ref.Checkout != "" && ref.Checkout != p.Checkout {
		return project.Project{}, app.Fail("invalid_path", "record is outside this conversation's project scope")
	}
	return p, nil
}
func (s *Server) documentCall(name string, args json.RawMessage) (any, error) {
	if name == "factory_projects_list" {
		var a struct{}
		if e := strict(args, &a); e != nil {
			return nil, e
		}
		return map[string]any{"projects": s.config.Projects}, nil
	}
	p, e := s.selected(knowledge.Ref{})
	if e != nil {
		return nil, e
	}
	switch name {
	case "factory_nodes_list":
		var a struct {
			Type  string `json:"type,omitempty"`
			Query string `json:"query,omitempty"`
			Limit int    `json:"limit,omitempty"`
		}
		if e = strict(args, &a); e != nil {
			return nil, e
		}
		if a.Limit == 0 {
			a.Limit = 50
		}
		if a.Limit < 1 || a.Limit > 200 {
			return nil, app.Fail("invalid_arguments", "limit must be 1 to 200")
		}
		snap, e := app.Load(s.config.Projects)
		if e != nil {
			return nil, e
		}
		b := snap.Bundle(p.Checkout)
		if b == nil {
			return nil, app.Fail("not_found", "project unavailable")
		}
		docs := []*knowledge.Document{}
		for _, d := range b.Docs {
			if a.Type != "" && d.Type != a.Type || a.Query != "" && !strings.Contains(strings.ToLower(d.Title+" "+d.Ref.Path), strings.ToLower(a.Query)) {
				continue
			}
			docs = append(docs, d)
		}
		sort.Slice(docs, func(i, j int) bool { return docs[i].Ref.Path < docs[j].Ref.Path })
		total := len(docs)
		if len(docs) > a.Limit {
			docs = docs[:a.Limit]
		}
		return map[string]any{"nodes": docs, "total": total, "truncated": len(docs) < total, "snapshot_revision": snap.Revision}, nil
	case "factory_definitions_list":
		var a struct{}
		if e = strict(args, &a); e != nil {
			return nil, e
		}
		snap, e := app.Load(s.config.Projects)
		if e != nil {
			return nil, e
		}
		b := snap.Bundle(p.Checkout)
		if b == nil {
			return nil, app.Fail("not_found", "project unavailable")
		}
		return map[string]any{"definitions": b.Definitions}, nil
	case "factory_validate":
		var a struct {
			Level string `json:"level,omitempty"`
		}
		if e = strict(args, &a); e != nil {
			return nil, e
		}
		if a.Level == "" {
			a.Level = "factory"
		}
		if a.Level != "okf" && a.Level != "factory" {
			return nil, app.Fail("invalid_arguments", "level must be okf or factory")
		}
		snap, e := app.Load(s.config.Projects)
		if e != nil {
			return nil, e
		}
		return snap.Validate(a.Level), nil
	case "factory_node_create":
		var a struct {
			Type   string         `json:"type"`
			Slug   string         `json:"slug"`
			Inputs map[string]any `json:"inputs,omitempty"`
		}
		if e = strict(args, &a); e != nil {
			return nil, e
		}
		return app.Create(s.config.Projects, p, a.Type, a.Slug, a.Inputs)
	case "factory_node_read", "factory_node_links", "factory_node_update", "factory_node_delete":
		var a nodeArgs
		if e = strict(args, &a); e != nil {
			return nil, e
		}
		p, e = s.selected(a.Ref)
		if e != nil {
			return nil, e
		}
		if a.Ref.Path == "" {
			return nil, app.Fail("invalid_arguments", "ref.path is required")
		}
		if name == "factory_node_update" {
			if (a.Set == nil) == (a.Document == nil) {
				return nil, app.Fail("invalid_arguments", "supply either set or document")
			}
			m := app.Mutation{Kind: "update", Path: a.Ref.Path, ExpectedRevision: a.ExpectedRevision, Set: a.Set}
			if a.Document != nil {
				m.Kind = "replace"
				m.Document = []byte(*a.Document)
			}
			return app.Mutate(s.config.Projects, p, m)
		}
		if name == "factory_node_delete" {
			if a.Set != nil || a.Document != nil {
				return nil, app.Fail("invalid_arguments", "delete does not accept set or document")
			}
			_, e := app.Mutate(s.config.Projects, p, app.Mutation{Kind: "delete", Path: a.Ref.Path, ExpectedRevision: a.ExpectedRevision})
			if e != nil {
				return nil, e
			}
			return map[string]any{"removed": a.Ref}, nil
		}
		if a.Set != nil || a.Document != nil || a.ExpectedRevision != "" {
			return nil, app.Fail("invalid_arguments", "read operation does not accept mutation fields")
		}
		snap, e := app.Load(s.config.Projects)
		if e != nil {
			return nil, e
		}
		b := snap.Bundle(p.Checkout)
		if b == nil || b.Docs[a.Ref.Path] == nil {
			return nil, app.Fail("not_found", "record not found in selected project")
		}
		d := b.Docs[a.Ref.Path]
		if name == "factory_node_links" {
			links := []app.Link{}
			for _, l := range snap.Links {
				if l.Source.Checkout == p.Checkout && l.Source.Path == a.Ref.Path || l.Resolved != nil && l.Resolved.Checkout == p.Checkout && l.Resolved.Path == a.Ref.Path {
					links = append(links, l)
				}
			}
			return map[string]any{"links": links}, nil
		}
		return map[string]any{"node": d, "source": string(d.Raw)}, nil
	default:
		return nil, app.Fail("unsupported_capability", "unknown tool")
	}
}

// Run serves MCP 2025-11-25 over newline-delimited JSON-RPC. EOF cancels pending waits.
func Run(ctx context.Context, in io.Reader, out io.Writer, config Config) error {
	s, e := New(config)
	if e != nil {
		return e
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 4096), 8<<20)
	var writeMu sync.Mutex
	var writeErr error
	var workers sync.WaitGroup
	var stateMu sync.Mutex
	pending := map[string]context.CancelFunc{}
	initialized := false
	sem := make(chan struct{}, 16)
	write := func(id json.RawMessage, result any, code int, message string) {
		writeMu.Lock()
		defer writeMu.Unlock()
		r := map[string]any{"jsonrpc": "2.0", "id": id}
		if code != 0 {
			r["error"] = map[string]any{"code": code, "message": message}
		} else {
			r["result"] = result
		}
		if writeErr == nil {
			writeErr = json.NewEncoder(out).Encode(r)
			if writeErr != nil {
				cancel()
			}
		}
	}
	for scanner.Scan() {
		raw := append([]byte(nil), scanner.Bytes()...)
		if len(bytes.TrimSpace(raw)) == 0 {
			continue
		}
		var req struct {
			JSONRPC string          `json:"jsonrpc"`
			ID      json.RawMessage `json:"id"`
			Method  string          `json:"method"`
			Params  json.RawMessage `json:"params"`
		}
		if json.Unmarshal(raw, &req) != nil {
			write(json.RawMessage("null"), nil, -32700, "invalid JSON")
			continue
		}
		if req.JSONRPC != "2.0" || req.Method == "" {
			write(req.ID, nil, -32600, "invalid request")
			continue
		}
		if len(req.ID) == 0 {
			if req.Method == "notifications/cancelled" {
				var p struct {
					RequestID json.RawMessage `json:"requestId"`
				}
				if json.Unmarshal(req.Params, &p) == nil {
					stateMu.Lock()
					if c := pending[string(p.RequestID)]; c != nil {
						c()
					}
					stateMu.Unlock()
				}
			}
			continue
		}
		if req.Method == "initialize" {
			if initialized {
				write(req.ID, nil, -32600, "already initialized")
				continue
			}
			initialized = true
			write(req.ID, map[string]any{"protocolVersion": "2025-11-25", "capabilities": map[string]any{"tools": map[string]any{}, "resources": map[string]any{}}, "serverInfo": map[string]any{"name": "factoryctl", "version": "0.1.0-dev"}, "instructions": "Document tools are project scoped. User answers come only from the attached trusted surface. Never synthesize answers; wait returns pending on timeout. Generated views require trusted host permission and have a plain-text fallback. External MCP Apps interoperability is unqualified."}, 0, "")
			continue
		}
		if !initialized {
			write(req.ID, nil, -32002, "initialize first")
			continue
		}
		switch req.Method {
		case "ping":
			write(req.ID, map[string]any{}, 0, "")
		case "tools/list":
			write(req.ID, map[string]any{"tools": Tools()}, 0, "")
		case "resources/list":
			write(req.ID, map[string]any{"resources": []any{}}, 0, "")
		case "resources/templates/list":
			write(req.ID, map[string]any{"resourceTemplates": []any{}}, 0, "")
		case "tools/call":
			var call struct {
				Request
				Meta map[string]json.RawMessage `json:"_meta,omitempty"`
			}
			if e := strict(req.Params, &call); e != nil {
				write(req.ID, nil, -32602, e.Error())
				continue
			}
			select {
			case sem <- struct{}{}:
			default:
				write(req.ID, nil, -32000, "too many concurrent requests")
				continue
			}
			stateMu.Lock()
			if _, exists := pending[string(req.ID)]; exists {
				stateMu.Unlock()
				<-sem
				write(req.ID, nil, -32600, "duplicate active request ID")
				continue
			}
			callctx, cc := context.WithCancel(ctx)
			pending[string(req.ID)] = cc
			stateMu.Unlock()
			workers.Add(1)
			go func(id json.RawMessage, c Request) {
				defer workers.Done()
				defer func() { stateMu.Lock(); delete(pending, string(id)); stateMu.Unlock(); cc(); <-sem }()
				result, e := s.Call(callctx, c.Name, c.Arguments)
				r := map[string]any{}
				if e != nil {
					var ae *app.Error
					if !errors.As(e, &ae) {
						ae = &app.Error{Code: "io_error", Message: e.Error()}
					}
					r["isError"] = true
					r["structuredContent"] = map[string]any{"error": ae}
					r["content"] = []any{map[string]any{"type": "text", "text": ae.Error()}}
				} else {
					data, err := json.Marshal(result)
					if err != nil {
						r["isError"] = true
						r["content"] = []any{map[string]any{"type": "text", "text": "io_error: cannot encode result"}}
					} else {
						r["content"] = []any{map[string]any{"type": "text", "text": string(data)}}
						r["structuredContent"] = result
					}
				}
				write(id, r, 0, "")
			}(append(json.RawMessage(nil), req.ID...), call.Request)
		default:
			write(req.ID, nil, -32601, "method not supported")
		}
	}
	cancel()
	workers.Wait()
	if writeErr != nil {
		return writeErr
	}
	return scanner.Err()
}
