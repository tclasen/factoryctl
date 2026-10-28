// Package web presents a read-only, loopback-only view of a document snapshot.
package web

import (
	"bytes"
	"context"
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/tclasen/factoryctl/internal/app"
	"github.com/tclasen/factoryctl/internal/project"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

//go:embed assets/*
var assets embed.FS

type Source struct {
	Label    string `json:"label"`
	Resource string `json:"resource"`
	ID       string `json:"id,omitempty"`
}
type Document struct {
	ID                  string   `json:"id"`
	Path                string   `json:"path"`
	Title               string   `json:"title"`
	Type                string   `json:"type"`
	Project             string   `json:"project"`
	IntentState         string   `json:"intent_state"`
	IntentBasis         string   `json:"intent_basis"`
	IntentOwner         string   `json:"intent_owner"`
	IntentSource        string   `json:"intent_source"`
	IntentAuthorization string   `json:"intent_authorization"`
	Revision            string   `json:"revision"`
	Sources             []Source `json:"sources"`
	Links               []Source `json:"links"`
	HTML                string   `json:"rendered_html,omitempty"`
	Raw                 string   `json:"raw,omitempty"`
}

type Handler struct {
	projects    []project.Project
	chat        *conversationHost
	documents   []Document
	details     map[string]Document
	diagnostics int
}

func value(m map[string]any, k string) string { s, _ := m[k].(string); return s }
func documentID(checkout, p string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(checkout + "\x00" + p))
}

// IntentState exposes only explicit metadata and explicit proposal defaults.
// It never promotes a Git author, stable document, passing test or source to acceptance.
func IntentState(m map[string]any) (string, string) {
	state := value(m, "intent_state")
	switch state {
	case "requested", "proposed", "accepted", "superseded", "unresolved":
		return state, "Explicit intent_state in source"
	}
	if state != "" {
		return "unresolved", "Unrecognized intent_state; inspect source"
	}
	factory, _ := m["factory"].(map[string]any)
	assertion, _ := factory["assertion"].(map[string]any)
	if value(assertion, "certainty") == "proposal" {
		return "proposed", "Document explicitly labels its assertion as a proposal"
	}
	return "unresolved", "No explicit intent state; acceptance is not inferred"
}

func New(projects []project.Project) (*Handler, error) {
	snap, err := app.Load(projects)
	if err != nil {
		return nil, err
	}
	h := &Handler{projects: projects, documents: []Document{}, details: map[string]Document{}, diagnostics: len(snap.Diagnostics)}
	byPath := map[string]string{}
	for _, b := range snap.Bundles {
		for p := range b.Docs {
			byPath[b.Project.Checkout+"\x00"+p] = documentID(b.Project.Checkout, p)
		}
	}
	resolve := func(checkout, from, target string) string {
		u, e := url.Parse(target)
		if e != nil || u.IsAbs() || u.Host != "" || u.RawQuery != "" {
			return ""
		}
		dest := u.Path
		if dest == "" {
			dest = from
		} else if strings.HasPrefix(dest, "/") {
			dest = strings.TrimPrefix(dest, "/")
		} else {
			dest = path.Join(path.Dir(from), dest)
		}
		if dest == ".." || strings.HasPrefix(dest, "../") || strings.Contains(dest, "\\") {
			return ""
		}
		return byPath[checkout+"\x00"+dest]
	}
	md := goldmark.New(goldmark.WithExtensions(extension.Table, extension.Strikethrough), goldmark.WithParserOptions(parser.WithAutoHeadingID()))
	for _, b := range snap.Bundles {
		for p, d := range b.Docs {
			state, basis := IntentState(d.Frontmatter)
			doc := Document{ID: documentID(b.Project.Checkout, p), Path: p, Title: d.Title, Type: d.Type, Project: b.Project.ID, IntentState: state, IntentBasis: basis, IntentOwner: value(d.Frontmatter, "intent_owner"), IntentSource: value(d.Frontmatter, "intent_source"), IntentAuthorization: value(d.Frontmatter, "intent_authorization"), Revision: d.Revision, Sources: []Source{}, Links: []Source{}, Raw: string(d.Raw)}
			if doc.Title == "" {
				doc.Title = p
			}
			source := []byte(d.Body)
			node := md.Parser().Parse(text.NewReader(source))
			// Remove images entirely: even a harmless image can leak local reading activity.
			var images []ast.Node
			_ = ast.Walk(node, func(n ast.Node, enter bool) (ast.WalkStatus, error) {
				if !enter {
					return ast.WalkContinue, nil
				}
				switch n := n.(type) {
				case *ast.Image:
					images = append(images, n)
					return ast.WalkSkipChildren, nil
				case *ast.Link:
					target := string(n.Destination)
					id := resolve(b.Project.Checkout, p, target)
					doc.Links = append(doc.Links, Source{Label: target, Resource: target, ID: id})
					if id != "" {
						u, _ := url.Parse(target)
						n.Destination = []byte("#doc=" + url.QueryEscape(id) + "&anchor=" + url.QueryEscape(u.Fragment))
					} else {
						n.Destination = []byte("#unavailable")
					}
				case *ast.AutoLink:
					// External autolinks are inert text. No navigation or resource fetching.
					n.SetAttributeString("data-external", []byte("true"))
				}
				return ast.WalkContinue, nil
			})
			for _, n := range images {
				n.Parent().ReplaceChild(n.Parent(), n, ast.NewString([]byte("[Image omitted in local document view]")))
			}
			// Convert autolinks to inert text before rendering.
			var auto []*ast.AutoLink
			_ = ast.Walk(node, func(n ast.Node, e bool) (ast.WalkStatus, error) {
				if e {
					if a, ok := n.(*ast.AutoLink); ok {
						auto = append(auto, a)
					}
				}
				return ast.WalkContinue, nil
			})
			for _, a := range auto {
				a.Parent().ReplaceChild(a.Parent(), a, ast.NewString(a.Label(source)))
			}
			var buf bytes.Buffer
			if err := md.Renderer().Render(&buf, source, node); err != nil {
				return nil, err
			}
			doc.HTML = buf.String()
			if doc.IntentSource != "" {
				doc.Sources = append(doc.Sources, Source{Label: "Intent source", Resource: doc.IntentSource, ID: resolve(b.Project.Checkout, p, doc.IntentSource)})
			}
			if sources, ok := d.Frontmatter["sources"].([]any); ok {
				for _, s := range sources {
					m, ok := s.(map[string]any)
					if !ok {
						continue
					}
					resource := value(m, "resource")
					if resource == "" {
						continue
					}
					label := value(m, "title")
					if label == "" {
						label = value(m, "id")
					}
					doc.Sources = append(doc.Sources, Source{Label: label, Resource: resource, ID: resolve(b.Project.Checkout, p, resource)})
				}
			}
			for _, field := range []string{"supersedes", "superseded_by"} {
				if target := value(d.Frontmatter, field); target != "" {
					doc.Links = append(doc.Links, Source{Label: field, Resource: target, ID: resolve(b.Project.Checkout, p, target)})
				}
			}
			h.details[doc.ID] = doc
			doc.HTML = ""
			doc.Raw = ""
			h.documents = append(h.documents, doc)
		}
	}
	sort.Slice(h.documents, func(i, j int) bool {
		if h.documents[i].Project != h.documents[j].Project {
			return h.documents[i].Project < h.documents[j].Project
		}
		return h.documents[i].Path < h.documents[j].Path
	})
	return h, nil
}

func localHost(hostport string) bool {
	host := hostport
	if strings.Contains(hostport, ":") {
		var err error
		host, _, err = net.SplitHostPort(hostport)
		if err != nil {
			return false
		}
	}
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self'; connect-src 'self'; img-src 'none'; frame-src 'none'; frame-ancestors 'none'; base-uri 'none'; form-action 'none'")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Cache-Control", "no-store")
	if !localHost(r.Host) {
		http.Error(w, "Loopback host required", http.StatusForbidden)
		return
	}
	if origin := r.Header.Get("Origin"); origin != "" {
		u, e := url.Parse(origin)
		if e != nil || u.Scheme != "http" || u.Host != r.Host {
			http.Error(w, "Same-origin access required", http.StatusForbidden)
			return
		}
	}
	if r.Header.Get("Sec-Fetch-Site") == "cross-site" {
		http.Error(w, "Cross-site access refused", http.StatusForbidden)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/view/") {
		h.viewHTTP(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/api/conversations") || r.URL.Path == "/internal/tool" || r.URL.Path == "/api/chat-settings" {
		h.conversationHTTP(w, r)
		return
	}
	if h.chat != nil {
		w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self'; connect-src 'self'; img-src 'none'; frame-src 'self'; frame-ancestors 'none'; base-uri 'none'; form-action 'none'")
	}
	if r.Method != "GET" && r.Method != "HEAD" {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "Read-only document view", http.StatusMethodNotAllowed)
		return
	}
	writeJSON := func(v any) {
		data, e := json.Marshal(v)
		if e != nil {
			http.Error(w, "Cannot encode snapshot", 500)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if r.Method != "HEAD" {
			_, _ = w.Write(data)
		}
	}
	switch r.URL.Path {
	case "/api/documents":
		writeJSON(map[string]any{"documents": h.documents, "diagnostics": h.diagnostics, "read_only": true, "snapshot": "Loaded at launch; restart to refresh"})
	case "/api/document":
		d, ok := h.details[r.URL.Query().Get("id")]
		if !ok {
			http.NotFound(w, r)
			return
		}
		writeJSON(d)
	case "/", "/app.js", "/style.css":
		name := strings.TrimPrefix(r.URL.Path, "/")
		if name == "" {
			name = "index.html"
		}
		b, e := assets.ReadFile("assets/" + name)
		if e != nil {
			http.NotFound(w, r)
			return
		}
		content := map[string]string{"index.html": "text/html; charset=utf-8", "app.js": "text/javascript; charset=utf-8", "style.css": "text/css; charset=utf-8"}
		w.Header().Set("Content-Type", content[name])
		if r.Method != "HEAD" {
			_, _ = w.Write(b)
		}
	default:
		http.NotFound(w, r)
	}
}

// Listen never binds a wildcard, external address or a DNS-resolved hostname.
func Listen(address string) (net.Listener, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, fmt.Errorf("listen address must be a loopback IP and port: %w", err)
	}
	if host == "localhost" {
		host = "127.0.0.1"
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		return nil, fmt.Errorf("web view must listen on a loopback IP")
	}
	return net.Listen("tcp", net.JoinHostPort(host, port))
}
func Serve(projects []project.Project, address string, out io.Writer) error {
	return ServeWithRuntime(projects, address, out, nil)
}
func ServeWithRuntime(projects []project.Project, address string, out io.Writer, cfg *RuntimeConfig) error {
	h, e := New(projects)
	if e != nil {
		return e
	}
	listener, e := Listen(address)
	if e != nil {
		return e
	}
	defer listener.Close()
	defer h.Close()
	if cfg != nil {
		h.EnableConversations(*cfg)
		h.chat.base = "http://" + listener.Addr().String()
	}
	if _, e = fmt.Fprintf(out, "factoryctl document explorer: http://%s\nRead-only local snapshot; Ctrl-C stops the server.\n", listener.Addr()); e != nil {
		return e
	}
	if h.chat != nil {
		if _, e = fmt.Fprintf(out, "Conversation access: http://%s/#host=%s\n", listener.Addr(), h.chat.secret); e != nil {
			return e
		}
	}
	server := &http.Server{Handler: h, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 45 * time.Second, IdleTimeout: 30 * time.Second, MaxHeaderBytes: 16 << 10}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() { <-ctx.Done(); h.Close(); _ = server.Close() }()
	err := server.Serve(listener)
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}
