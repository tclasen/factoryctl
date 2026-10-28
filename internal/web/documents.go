package web

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/tclasen/factoryctl/internal/app"
	"github.com/tclasen/factoryctl/internal/knowledge"
	"github.com/tclasen/factoryctl/internal/project"
)

func snapshotIdentity(s *app.Snapshot) string {
	b, _ := json.Marshal(s.Diagnostics)
	return knowledge.Hash(append([]byte(s.Revision), b...))
}
func (h *Handler) refresh() error {
	h.docMu.Lock()
	defer h.docMu.Unlock()
	snap, e := app.LoadCached(h.projects, h.cached)
	if e != nil {
		return e
	}
	if snapshotIdentity(snap) == h.snapshot {
		h.cached = snap
		return nil
	}
	next, e := renderHandler(h.projects, snap)
	if e != nil {
		return e
	}
	h.documents = next.documents
	h.details = next.details
	h.diagnostics = next.diagnostics
	h.snapshot = next.snapshot
	h.cached = snap
	return nil
}
func documentJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(v)
}
func documentError(w http.ResponseWriter, e error) {
	var a *app.Error
	if !errors.As(e, &a) {
		a = &app.Error{Code: "io_error", Message: e.Error()}
	}
	status := 422
	switch a.Code {
	case "not_found":
		status = 404
	case "invalid_path", "invalid_arguments":
		status = 400
	case "stale_definition", "stale_revision", "schema_conflict", "ambiguous_project", "referenced_node", "already_exists", "cursor_expired":
		status = 409
	case "io_error":
		status = 500
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	documentJSON(w, map[string]any{"error": a})
}
func decodeDocument(r *http.Request, v any) error {
	d := json.NewDecoder(io.LimitReader(r.Body, knowledge.MaxDocument+65536))
	d.DisallowUnknownFields()
	if e := d.Decode(v); e != nil {
		return app.Fail("invalid_arguments", e.Error())
	}
	var extra any
	if e := d.Decode(&extra); e != io.EOF {
		return app.Fail("invalid_arguments", "one JSON object required")
	}
	return nil
}
func (h *Handler) selected(ref knowledge.Ref) (project.Project, error) {
	var found []project.Project
	for _, p := range h.projects {
		if (ref.Checkout == "" || p.Checkout == ref.Checkout) && (ref.ProjectID == "" || p.ID == ref.ProjectID) {
			found = append(found, p)
		}
	}
	if len(found) != 1 {
		return project.Project{}, app.Fail("ambiguous_project", "select one loaded project instance")
	}
	return found[0], nil
}
func queryRef(r *http.Request) knowledge.Ref {
	return knowledge.Ref{ProjectID: r.URL.Query().Get("project_id"), Checkout: r.URL.Query().Get("checkout"), Path: r.URL.Query().Get("path")}
}
func (h *Handler) documentHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" && subtle.ConstantTimeCompare([]byte(r.Header.Get("X-Factory-Nonce")), []byte(h.nonce)) != 1 {
		w.WriteHeader(403)
		documentJSON(w, map[string]any{"error": map[string]any{"code": "forbidden", "message": "Session nonce required"}})
		return
	}
	if r.URL.Path == "/api/v1/session" && r.Method == "GET" {
		documentJSON(w, map[string]any{"nonce": h.nonce})
		return
	}
	if r.URL.Path == "/api/v1/delivery" && r.Method == "GET" {
		p, e := h.selected(queryRef(r))
		if e != nil {
			documentError(w, e)
			return
		}
		h.docMu.RLock()
		d, ok := h.details[documentID(p.Checkout, r.URL.Query().Get("path"))]
		h.docMu.RUnlock()
		if !ok {
			documentError(w, app.Fail("not_found", "record missing"))
			return
		}
		documentJSON(w, h.documentDelivery(d))
		return
	}
	if strings.HasPrefix(r.URL.Path, "/api/v1/custom") {
		h.customHTTP(w, r)
		return
	}
	if e := h.refresh(); e != nil {
		documentError(w, e)
		return
	}
	if r.Method == "GET" || r.Method == "HEAD" {
		h.docMu.RLock()
		defer h.docMu.RUnlock()
		switch r.URL.Path {
		case "/api/v1/projects":
			documentJSON(w, map[string]any{"projects": h.projects, "diagnostics": h.diagnostics})
			return
		case "/api/v1/nodes":
			q := r.URL.Query()
			rows := []Document{}
			for _, d := range h.documents {
				if q.Get("checkout") != "" && d.Ref.Checkout != q.Get("checkout") {
					continue
				}
				if q.Get("project_id") != "" && d.Project != q.Get("project_id") {
					continue
				}
				if q.Get("type") != "" && d.Type != q.Get("type") {
					continue
				}
				if !strings.Contains(strings.ToLower(d.Title+" "+d.Path), strings.ToLower(q.Get("query"))) {
					continue
				}
				rows = append(rows, d)
			}
			scopeBytes, _ := json.Marshal([]string{q.Get("checkout"), q.Get("project_id"), q.Get("type"), q.Get("query")})
			scopeHash := knowledge.Hash(scopeBytes)
			offset := 0
			if c := q.Get("cursor"); c != "" {
				bits := strings.Split(c, ":")
				if len(bits) != 3 || bits[0] != h.snapshot || bits[1] != scopeHash {
					documentError(w, app.Fail("cursor_expired", "snapshot or filters changed; restart list"))
					return
				}
				var parseErr error
				offset, parseErr = strconv.Atoi(bits[2])
				if parseErr != nil || offset < 0 {
					documentError(w, app.Fail("invalid_arguments", "invalid cursor offset"))
					return
				}
			}
			limit, _ := strconv.Atoi(q.Get("limit"))
			if limit < 1 || limit > 200 {
				limit = 50
			}
			if offset < 0 || offset > len(rows) {
				documentError(w, app.Fail("invalid_arguments", "invalid cursor"))
				return
			}
			end := offset + limit
			if end > len(rows) {
				end = len(rows)
			}
			next := ""
			if end < len(rows) {
				next = h.snapshot + ":" + scopeHash + ":" + strconv.Itoa(end)
			}
			documentJSON(w, map[string]any{"nodes": rows[offset:end], "next_cursor": next, "snapshot": h.snapshot})
			return
		case "/api/v1/node", "/api/v1/links":
			ref := queryRef(r)
			p, e := h.selected(ref)
			if e != nil {
				documentError(w, e)
				return
			}
			d, ok := h.details[documentID(p.Checkout, ref.Path)]
			if !ok {
				documentError(w, app.Fail("not_found", "record not found"))
				return
			}
			if r.URL.Path == "/api/v1/links" {
				documentJSON(w, map[string]any{"incoming": d.Backlinks, "outgoing": d.Links})
			} else {
				if d.Frontmatter["delivery"] != nil || d.Frontmatter["delivery_ref"] != nil {
					d.Delivery = h.documentDelivery(d)
				}
				documentJSON(w, d)
			}
			return
		case "/api/v1/definitions", "/api/v1/validate", "/api/v1/actions", "/api/v1/extensions":
			p, e := h.selected(queryRef(r))
			if e != nil {
				documentError(w, e)
				return
			}
			snap, e := app.Load(h.projects)
			if e != nil {
				documentError(w, e)
				return
			}
			switch r.URL.Path {
			case "/api/v1/definitions":
				defs := map[string]any{}
				for typ, def := range snap.Bundle(p.Checkout).Definitions {
					defs[typ] = map[string]any{"type": typ, "fields": def.Fields, "raw": def.Raw, "enabled": def.Enabled, "invalid": def.Invalid, "definition_hash": app.DefinitionHash(def)}
				}
				documentJSON(w, defs)
			case "/api/v1/validate":
				documentJSON(w, snap.Validate("factory"))
			case "/api/v1/actions":
				a, e := app.Actions(h.projects, p, r.URL.Query().Get("path"))
				if e != nil {
					documentError(w, e)
					return
				}
				documentJSON(w, a)
			case "/api/v1/extensions":
				a, e := app.ExtensionCatalog(h.projects, p)
				if e != nil {
					documentError(w, e)
					return
				}
				documentJSON(w, a)
			}
			return
		}
	}
	if r.URL.Path == "/api/v1/init" && r.Method == "POST" {
		var v struct {
			Ref       knowledge.Ref `json:"ref"`
			ProjectID string        `json:"project_id"`
			Starter   string        `json:"starter"`
		}
		if e := decodeDocument(r, &v); e != nil {
			documentError(w, e)
			return
		}
		p, e := h.selected(v.Ref)
		if e != nil {
			documentError(w, e)
			return
		}
		out, e := app.Init(p.Root, v.ProjectID, v.Starter)
		if e != nil {
			documentError(w, e)
			return
		}
		_ = h.refresh()
		documentJSON(w, out)
		return
	}
	if r.URL.Path == "/api/v1/actions/render" && r.Method == "POST" {
		var v struct {
			Ref    knowledge.Ref  `json:"ref"`
			Prompt string         `json:"prompt"`
			Inputs map[string]any `json:"inputs"`
		}
		if e := decodeDocument(r, &v); e != nil {
			documentError(w, e)
			return
		}
		p, e := h.selected(v.Ref)
		if e != nil {
			documentError(w, e)
			return
		}
		out, e := app.RenderPrompt(h.projects, p, v.Prompt, v.Ref.Path, v.Inputs)
		if e != nil {
			documentError(w, e)
			return
		}
		documentJSON(w, out)
		return
	}
	if r.URL.Path == "/api/v1/extensions/call" && r.Method == "POST" {
		var v struct {
			Ref       knowledge.Ref  `json:"ref"`
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		}
		if e := decodeDocument(r, &v); e != nil {
			documentError(w, e)
			return
		}
		p, e := h.selected(v.Ref)
		if e != nil {
			documentError(w, e)
			return
		}
		out, e := app.CallExtension(h.projects, p, v.Name, v.Arguments)
		if e != nil {
			documentError(w, e)
			return
		}
		documentJSON(w, out)
		return
	}
	if (r.URL.Path == "/api/v1/node" && (r.Method == "PUT" || r.Method == "DELETE")) || (r.URL.Path == "/api/v1/nodes" && r.Method == "POST") {
		var v struct {
			Ref                knowledge.Ref  `json:"ref"`
			ExpectedRevision   string         `json:"expected_revision"`
			ExpectedDefinition string         `json:"expected_definition"`
			Set                map[string]any `json:"set"`
			Document           *string        `json:"document"`
			Type               string         `json:"type"`
			Slug               string         `json:"slug"`
			Inputs             map[string]any `json:"inputs"`
		}
		if e := decodeDocument(r, &v); e != nil {
			documentError(w, e)
			return
		}
		p, e := h.selected(v.Ref)
		if e != nil {
			documentError(w, e)
			return
		}
		if v.ExpectedDefinition == "" {
			documentError(w, app.Fail("invalid_arguments", "expected_definition required"))
			return
		}
		var d *knowledge.Document
		if r.Method == "POST" {
			d, e = app.CreateChecked(h.projects, p, v.Type, v.Slug, v.Inputs, v.ExpectedDefinition)
		} else {
			m := app.Mutation{Path: v.Ref.Path, ExpectedRevision: v.ExpectedRevision, ExpectedDefinition: v.ExpectedDefinition}
			if r.Method == "DELETE" {
				m.Kind = "delete"
			} else if v.Document != nil && v.Set == nil {
				m.Kind = "replace"
				m.Document = []byte(*v.Document)
			} else if v.Document == nil && v.Set != nil {
				m.Kind = "update"
				m.Set = v.Set
			} else {
				documentError(w, app.Fail("invalid_arguments", "supply exactly one of set or document"))
				return
			}
			d, e = app.Mutate(h.projects, p, m)
		}
		if e != nil {
			documentError(w, e)
			return
		}
		if e = h.refresh(); e != nil {
			documentError(w, e)
			return
		}
		if r.Method == "DELETE" || d == nil {
			documentJSON(w, map[string]any{"removed": v.Ref})
			return
		}
		h.docMu.RLock()
		out := h.details[documentID(p.Checkout, d.Ref.Path)]
		h.docMu.RUnlock()
		documentJSON(w, out)
		return
	}
	w.Header().Set("Allow", "GET, POST, PUT, DELETE")
	w.WriteHeader(405)
	documentJSON(w, map[string]any{"error": fmt.Sprintf("unsupported method for %s", r.URL.Path)})
}
