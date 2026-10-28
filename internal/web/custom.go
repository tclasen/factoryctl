package web

import (
	"bytes"

	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/tclasen/factoryctl/internal/app"
	"github.com/tclasen/factoryctl/internal/knowledge"
	"github.com/tclasen/factoryctl/internal/project"
	"github.com/tclasen/factoryctl/internal/safeio"
)

type customGrant struct {
	Project              project.Project
	Ref                  knowledge.Ref
	Selection            map[string]bool
	Digest, Nonce, Entry string
	Files                map[string][]byte
	Definition           string
}

var scriptAssetTag = regexp.MustCompile(`(?is)<script\b[^>]*\bsrc\s*=\s*["']([^"']+)["'][^>]*>\s*</script\s*>`)
var styleAssetTag = regexp.MustCompile(`(?is)<link\b[^>]*\bhref\s*=\s*["']([^"']+)["'][^>]*>`)
var assetReferences = regexp.MustCompile(`(?i)(?:src\s*=\s*["']|href\s*=\s*["']|url\(\s*["']?|(?:import|from)\s*["'])([^"'\s)>]+)`)

func readAsset(root, name string) ([]byte, error) {
	if !strings.HasPrefix(name, ".factory/") || path.Clean(name) != name || strings.ContainsAny(name, "\\:%\x00?#") {
		return nil, app.Fail("invalid_path", "custom assets must be normalized .factory/ paths")
	}
	current := root
	for _, part := range strings.Split(name, "/") {
		current = filepath.Join(current, part)
		st, e := os.Lstat(current)
		if e != nil {
			return nil, e
		}
		if st.Mode()&os.ModeSymlink != 0 {
			return nil, app.Fail("invalid_path", "asset symlinks forbidden")
		}
	}
	scope, e := os.OpenRoot(root)
	if e != nil {
		return nil, e
	}
	defer scope.Close()
	b, e := safeio.ReadRoot(scope, filepath.FromSlash(name), 1<<20)
	if e == nil && len(b) > 1<<20 {
		return nil, app.Fail("invalid_document", "asset exceeds 1 MiB")
	}
	return b, e
}
func (h *Handler) customBundle(ref knowledge.Ref) (customGrant, error) {
	p, e := h.selected(ref)
	if e != nil {
		return customGrant{}, e
	}
	s, e := app.Load(h.projects)
	if e != nil {
		return customGrant{}, e
	}
	b := s.Bundle(p.Checkout)
	d := b.Docs[ref.Path]
	if d == nil {
		return customGrant{}, app.Fail("not_found", "record missing")
	}
	def := b.Definitions[d.Type]
	if !def.Enabled || def.Invalid || def.ExtensionInvalid {
		return customGrant{}, app.Fail("schema_conflict", "custom view definition is disabled or invalid; use generic view")
	}
	views := knowledge.Mapping(def.Raw["views"])
	entry := knowledge.String(views, "detail_template")
	if entry == "" {
		return customGrant{}, app.Fail("not_found", "no custom view declared; generic view available")
	}
	g := customGrant{Project: p, Ref: d.Ref, Selection: map[string]bool{}, Entry: entry, Files: map[string][]byte{}, Definition: app.DefinitionHash(def)}
	for _, candidate := range b.Docs {
		if candidate.Type == d.Type {
			g.Selection[documentID(p.Checkout, candidate.Ref.Path)] = true
		}
	}
	queue := []string{entry}
	for _, v := range knowledge.List(views["assets"]) {
		a, ok := v.(string)
		if !ok {
			return g, app.Fail("invalid_document", "view assets must be paths")
		}
		queue = append(queue, a)
	}
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		if _, ok := g.Files[name]; ok {
			continue
		}
		if len(g.Files) >= 64 {
			return g, app.Fail("invalid_document", "too many custom assets")
		}
		raw, e := readAsset(p.Root, name)
		if e != nil {
			return g, e
		}
		g.Files[name] = raw
		for _, m := range assetReferences.FindAllSubmatch(raw, -1) {
			target := string(m[1])
			if strings.HasPrefix(target, "#") || strings.HasPrefix(target, "data:") {
				continue
			}
			u, e := url.Parse(target)
			if e != nil || u.IsAbs() || u.Host != "" {
				return g, app.Fail("invalid_path", "external custom resources are disabled")
			}
			dest := target
			if !strings.HasPrefix(dest, ".factory/") {
				dest = path.Join(path.Dir(name), dest)
			}
			queue = append(queue, dest)
		}
	}
	keys := []string{}
	for k := range g.Files {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var digest bytes.Buffer
	digest.WriteString(g.Definition)
	for _, k := range keys {
		digest.WriteString(k)
		digest.Write(g.Files[k])
	}
	g.Digest = knowledge.Hash(digest.Bytes())
	return g, nil
}
func (h *Handler) validGrant(id string) (customGrant, error) {
	h.customMu.Lock()
	g, ok := h.grants[id]
	h.customMu.Unlock()
	if !ok {
		return g, app.Fail("not_found", "custom frame unavailable")
	}
	fresh, e := h.customBundle(g.Ref)
	if e != nil || fresh.Digest != g.Digest {
		h.customMu.Lock()
		delete(h.grants, id)
		h.customMu.Unlock()
		return g, app.Fail("stale_revision", "custom assets or definition changed; trust again")
	}
	return g, nil
}
func (h *Handler) customHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		g, e := h.customBundle(queryRef(r))
		if e != nil {
			documentError(w, e)
			return
		}
		documentJSON(w, map[string]any{"digest": g.Digest, "entry": g.Entry, "trust_required": true})
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var v struct {
		Ref              knowledge.Ref  `json:"ref"`
		Digest           string         `json:"digest"`
		Grant            string         `json:"grant"`
		Nonce            string         `json:"nonce"`
		Type             string         `json:"type"`
		ID               string         `json:"id"`
		ExpectedRevision string         `json:"expected_revision"`
		Set              map[string]any `json:"set"`
		Prompt           string         `json:"prompt"`
		Inputs           map[string]any `json:"inputs"`
	}
	if e := decodeDocument(r, &v); e != nil {
		documentError(w, e)
		return
	}
	if r.URL.Path == "/api/v1/custom/trust" {
		g, e := h.customBundle(v.Ref)
		if e != nil {
			documentError(w, e)
			return
		}
		if v.Digest != g.Digest {
			documentError(w, app.Fail("stale_revision", "assets changed; review and trust again"))
			return
		}
		id := randomToken()
		g.Nonce = randomToken()
		h.customMu.Lock()
		h.grants[id] = g
		h.customMu.Unlock()
		h.docMu.RLock()
		selection := []Document{}
		for key := range g.Selection {
			if d, ok := h.details[key]; ok {
				selection = append(selection, d)
			}
		}
		h.docMu.RUnlock()
		documentJSON(w, map[string]any{"selection": selection, "grant": id, "nonce": g.Nonce, "digest": g.Digest, "frame_url": "/custom/" + id + "/" + g.Entry})
		return
	}
	g, e := h.validGrant(v.Grant)
	if e != nil {
		documentError(w, e)
		return
	}
	if v.Nonce != g.Nonce || v.Digest != g.Digest || !g.Selection[v.ID] {
		documentError(w, app.Fail("invalid_arguments", "custom message is outside granted selection"))
		return
	}
	var ref knowledge.Ref
	h.docMu.RLock()
	d, ok := h.details[v.ID]
	h.docMu.RUnlock()
	if !ok {
		documentError(w, app.Fail("not_found", "selected record disappeared"))
		return
	}
	ref = d.Ref
	switch v.Type {
	case "requestEdit":
		_, e = app.Mutate(h.projects, g.Project, app.Mutation{Kind: "update", Path: ref.Path, ExpectedRevision: v.ExpectedRevision, ExpectedDefinition: g.Definition, Set: v.Set})
		if e != nil {
			documentError(w, e)
			return
		}
		_ = h.refresh()
		h.docMu.RLock()
		d = h.details[v.ID]
		h.docMu.RUnlock()
		documentJSON(w, d)
	case "requestAction":
		actions, e := app.Actions(h.projects, g.Project, ref.Path)
		if e != nil {
			documentError(w, e)
			return
		}
		declared := false
		for _, a := range actions {
			if a["path"] == v.Prompt {
				declared = true
			}
		}
		if !declared {
			documentError(w, app.Fail("invalid_arguments", "prompt is not a declared action"))
			return
		}
		out, e := app.RenderPrompt(h.projects, g.Project, v.Prompt, ref.Path, v.Inputs)
		if e != nil {
			documentError(w, e)
			return
		}
		documentJSON(w, out)
	default:
		documentError(w, app.Fail("invalid_arguments", "unsupported custom request"))
	}
}
func (h *Handler) customFrame(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		return
	}
	parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/custom/"), "/", 2)
	if len(parts) != 2 {
		http.NotFound(w, r)
		return
	}
	g, e := h.validGrant(parts[0])
	if e != nil {
		documentError(w, e)
		return
	}
	name := parts[1]
	raw, ok := g.Files[name]
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'unsafe-inline' 'self'; style-src 'unsafe-inline' 'self'; connect-src 'none'; img-src data:; font-src 'none'; frame-src 'none'; worker-src 'none'; sandbox allow-scripts; frame-ancestors 'self'; base-uri 'none'; form-action 'none'")
	switch path.Ext(name) {
	case ".html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		var model any
		h.docMu.RLock()
		selected := []Document{}
		for id := range g.Selection {
			if d, ok := h.details[id]; ok {
				selected = append(selected, d)
			}
		}
		model = map[string]any{"Record": h.details[documentID(g.Ref.Checkout, g.Ref.Path)], "Selection": selected}
		h.docMu.RUnlock()
		t, e := template.New(name).Parse(string(raw))
		if e != nil {
			http.Error(w, "Invalid custom template", 422)
			return
		}
		var b bytes.Buffer
		if e = t.Execute(&b, model); e != nil {
			http.Error(w, "Cannot render custom template", 422)
			return
		}
		html := scriptAssetTag.ReplaceAllStringFunc(b.String(), func(tag string) string {
			m := scriptAssetTag.FindStringSubmatch(tag)
			key := m[1]
			if !strings.HasPrefix(key, ".factory/") {
				key = path.Join(path.Dir(name), key)
			}
			raw, ok := g.Files[key]
			if !ok {
				return ""
			}
			return "<script>" + string(raw) + "</script>"
		})
		html = styleAssetTag.ReplaceAllStringFunc(html, func(tag string) string {
			m := styleAssetTag.FindStringSubmatch(tag)
			key := m[1]
			if !strings.HasPrefix(key, ".factory/") {
				key = path.Join(path.Dir(name), key)
			}
			raw, ok := g.Files[key]
			if !ok || path.Ext(key) != ".css" {
				return ""
			}
			return "<style>" + string(raw) + "</style>"
		})
		_, _ = w.Write([]byte(html))
	case ".css":
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		_, _ = w.Write(raw)
	case ".js":
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		_, _ = w.Write(raw)
	default:
		http.Error(w, fmt.Sprintf("unsupported asset %s", name), 415)
	}
}
