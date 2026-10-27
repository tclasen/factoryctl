package app

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tclasen/factoryctl/internal/knowledge"
	"github.com/tclasen/factoryctl/internal/project"
	"github.com/tclasen/factoryctl/internal/safeio"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func SafePath(root, p string, create bool) error {
	if p == "" || path.IsAbs(p) || p != path.Clean(p) || strings.Contains(p, "\\") || strings.Contains(p, "%") || strings.Contains(p, ":") || strings.ContainsRune(p, 0) || p == ".." || strings.HasPrefix(p, "../") || !strings.HasSuffix(p, ".md") {
		return Fail("invalid_path", "expected a normalized bundle-relative .md path")
	}
	base := strings.ToLower(path.Base(p))
	if base == "index.md" || base == "log.md" {
		return Fail("invalid_path", "reserved filename")
	}
	// Inspect the root itself as well as every descendant. Root-scoped OS calls
	// additionally prevent escaping if a non-cooperating editor races this check.
	current := root
	for _, part := range append([]string{""}, strings.Split(p, "/")...) {
		if part != "" {
			current = filepath.Join(current, part)
		}
		st, e := os.Lstat(current)
		if os.IsNotExist(e) && create {
			continue
		}
		if e != nil {
			return e
		}
		if st.Mode()&os.ModeSymlink != 0 {
			return Fail("invalid_path", "symlink components are forbidden")
		}
		if !st.IsDir() && !st.Mode().IsRegular() {
			return Fail("invalid_path", "special files are forbidden")
		}
	}
	return nil
}

type Mutation struct {
	Kind, Path, ExpectedRevision string
	Document                     []byte
	Set                          map[string]any
}

func Mutate(projects []project.Project, selected project.Project, m Mutation) (*knowledge.Document, error) {
	rootPath := filepath.Join(selected.Root, "factory")
	if e := SafePath(rootPath, m.Path, m.Kind == "create"); e != nil {
		return nil, e
	}
	if m.Path == "control/write-coordination.md" {
		return nil, Fail("invalid_path", "the stable lock anchor cannot be mutated")
	}
	if e := SafePath(rootPath, "control/write-coordination.md", false); e != nil {
		return nil, Fail("io_error", "missing or unsafe write-coordination anchor; initialize the bundle first")
	}
	root, e := os.OpenRoot(rootPath)
	if e != nil {
		return nil, e
	}
	defer root.Close()
	lock, e := safeio.OpenRegularRoot(root, "control/write-coordination.md", os.O_RDWR)
	if e != nil {
		return nil, e
	}
	defer lock.Close()
	if e = lockFile(lock); e != nil {
		return nil, Fail("stale_revision", "bundle writer is busy; refresh and retry")
	}
	defer unlockFile(lock)
	lockInfo, e := lock.Stat()
	if e != nil {
		return nil, e
	}
	currentLock, e := root.Stat("control/write-coordination.md")
	if e != nil || !os.SameFile(lockInfo, currentLock) {
		return nil, Fail("stale_revision", "lock anchor changed; refresh")
	}
	before, e := Load(projects)
	if e != nil {
		return nil, e
	}
	bundle := before.Bundle(selected.Checkout)
	if bundle == nil {
		return nil, Fail("not_found", "selected bundle is unavailable")
	}
	old := bundle.Docs[m.Path]
	if m.Kind == "create" {
		if _, e := root.Lstat(m.Path); e == nil {
			return nil, Fail("already_exists", "record already exists")
		} else if !os.IsNotExist(e) {
			return nil, e
		}
	} else {
		if old == nil {
			return nil, Fail("not_found", "record does not exist")
		}
		if m.ExpectedRevision == "" || m.ExpectedRevision != old.Revision {
			return nil, &Error{Code: "stale_revision", Message: "expected revision differs from disk", CurrentRevision: old.Revision}
		}
	}
	data := m.Document
	if m.Kind == "update" {
		if old == nil {
			return nil, Fail("not_found", "record does not exist")
		}
		known := map[string]bool{"type": true, "title": true, "description": true, "status": true, "resource": true, "tags": true, "generated": true, "verified": true, "sources": true, "stale_after": true, "factory": true}
		for k := range bundle.Definitions[old.Type].Fields {
			known[k] = true
		}
		for k := range m.Set {
			if !known[k] {
				return nil, Fail("invalid_document", "field is not declared: "+k)
			}
		}
		data, e = old.Edit(m.Set)
		if e != nil {
			return nil, Fail("lossless_edit_unavailable", e.Error())
		}
	}
	prior := before.Validate("factory")
	previous := map[string]int{}
	for _, d := range prior.Diagnostics {
		if d.Severity == "error" {
			previous[diagnosticKey(d)]++
		}
	}
	var candidate *knowledge.Document
	if m.Kind == "delete" {
		delete(bundle.Docs, m.Path)
	} else {
		candidate = knowledge.Parse(m.Path, data)
		candidate.Ref = knowledge.Ref{ProjectID: selected.ID, Checkout: selected.Checkout, Path: m.Path}
		if old != nil && old.Type == "Factory" && (knowledge.String(old.Frontmatter, "project_id") != knowledge.String(candidate.Frontmatter, "project_id") || knowledge.String(old.Frontmatter, "bundle_uri") != knowledge.String(candidate.Frontmatter, "bundle_uri")) {
			return nil, Fail("unsupported_capability", "identity changes require an external edit followed by project rediscovery")
		}
		bundle.Docs[m.Path] = candidate
	}
	post := before.Validate("factory")
	rejected := []knowledge.Diagnostic{}
	for _, d := range post.Diagnostics {
		if d.Severity != "error" {
			continue
		}
		k := diagnosticKey(d)
		if previous[k] > 0 {
			previous[k]--
			continue
		}
		rejected = append(rejected, d)
	}
	// A changed document must itself be valid even when an old error is unchanged.
	for _, d := range post.Diagnostics {
		if d.Severity == "error" && d.Checkout == selected.Checkout && d.Path == m.Path {
			rejected = append(rejected, d)
		}
	}
	if len(rejected) > 0 {
		code := "invalid_document"
		if m.Kind == "delete" {
			code = "referenced_node"
		}
		return nil, &Error{Code: code, Message: "proposed document graph fails structural validation", Diagnostics: rejected}
	}
	// Recheck all observed files, including affected incoming records. This does
	// not create a transaction with external editors or other bundle writers.
	for _, b := range before.Bundles {
		for p, d := range b.Docs {
			if b.Project.Checkout == selected.Checkout && p == m.Path {
				continue
			}
			raw, e := ReadBounded(filepath.Join(b.Project.Root, "factory", filepath.FromSlash(p)))
			if e != nil || knowledge.Hash(raw) != d.Revision {
				return nil, Fail("stale_revision", "validation dependency changed; refresh")
			}
		}
	}
	if old != nil {
		raw, e := safeio.ReadRoot(root, m.Path, knowledge.MaxDocument)
		if e != nil || !bytes.Equal(raw, old.Raw) {
			return nil, Fail("stale_revision", "target changed during validation")
		}
	}
	if e := SafePath(rootPath, m.Path, m.Kind == "create"); e != nil {
		return nil, e
	}
	currentLock, e = root.Stat("control/write-coordination.md")
	if e != nil || !os.SameFile(lockInfo, currentLock) {
		return nil, Fail("stale_revision", "lock anchor changed during validation")
	}
	if m.Kind == "delete" {
		if e = root.Remove(m.Path); e != nil {
			return nil, e
		}
		return old, nil
	}
	if m.Kind == "create" {
		if e = root.MkdirAll(path.Dir(m.Path), 0755); e != nil {
			return nil, e
		}
		if e = SafePath(rootPath, m.Path, true); e != nil {
			return nil, e
		}
	}
	// Random temporary sibling, fsync before publication; exclusive hard-link
	// publication prevents create from replacing a file created by an outsider.
	random := make([]byte, 16)
	if _, e = rand.Read(random); e != nil {
		return nil, e
	}
	tmpName := path.Join(path.Dir(m.Path), ".factoryctl-"+hex.EncodeToString(random))
	tmp, e := root.OpenFile(tmpName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if e != nil {
		return nil, e
	}
	defer root.Remove(tmpName)
	mode := os.FileMode(0644)
	if old != nil {
		if st, e := root.Stat(m.Path); e == nil {
			mode = st.Mode().Perm()
		}
	}
	e = tmp.Chmod(mode)
	if e == nil {
		_, e = tmp.Write(data)
	}
	if e == nil {
		e = tmp.Sync()
	}
	ce := tmp.Close()
	if e == nil {
		e = ce
	}
	if e != nil {
		return nil, e
	}
	if m.Kind == "create" {
		e = root.Link(tmpName, m.Path)
		if os.IsExist(e) {
			return nil, Fail("already_exists", "record appeared during creation")
		}
	} else {
		e = root.Rename(tmpName, m.Path)
	}
	if e != nil {
		return nil, e
	}
	if dir, e := safeio.OpenDirectoryRoot(root, path.Dir(m.Path)); e == nil {
		_ = dir.Sync()
		_ = dir.Close()
	}
	return candidate, nil
}
func diagnosticKey(d knowledge.Diagnostic) string {
	return fmt.Sprintf("%s|%s|%s|%s|%s|%s", d.ProjectID+d.Checkout, d.Path, d.Code, d.Field, d.Target, d.Message)
}
