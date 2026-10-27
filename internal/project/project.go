// Package project discovers bounded collections using filesystem metadata only.
package project

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/tclasen/factoryctl/internal/safeio"
	"gopkg.in/yaml.v3"
)

type Project struct {
	ID             string `json:"project_id"`
	Name           string `json:"name"`
	Root           string `json:"root"`
	RepositoryRoot string `json:"repository_root"`
	Checkout       string `json:"checkout"`
	BundleURI      string `json:"bundle_uri,omitempty"`
}
type Diagnostic struct {
	Code      string `json:"code"`
	Severity  string `json:"severity"`
	ProjectID string `json:"project_id,omitempty"`
	Path      string `json:"path"`
	Message   string `json:"message"`
}
type Collection struct {
	Projects    []Project    `json:"projects"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

// FindRepository identifies the nearest enclosing working tree. Bare repositories
// are deliberately not projects. It never invokes Git or reads repository config.
func FindRepository(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	abs, err = filepath.EvalSymlinks(abs)
	if err != nil {
		return "", err
	}
	for p := abs; ; p = filepath.Dir(p) {
		marker := filepath.Join(p, ".git")
		st, e := os.Lstat(marker)
		if e == nil {
			if st.Mode()&os.ModeSymlink != 0 {
				return "", fmt.Errorf("symlink Git marker: %s", marker)
			}
			gitdir := marker
			if !st.IsDir() {
				if !st.Mode().IsRegular() {
					return "", fmt.Errorf("invalid Git marker: %s", marker)
				}
				b, e := readBounded(marker, 4096)
				if e != nil {
					return "", e
				}
				line := strings.TrimSpace(string(b))
				if !strings.HasPrefix(line, "gitdir: ") || strings.ContainsAny(line, "\r\n") {
					return "", fmt.Errorf("invalid worktree pointer: %s", marker)
				}
				gitdir = strings.TrimSpace(strings.TrimPrefix(line, "gitdir: "))
				if !filepath.IsAbs(gitdir) {
					gitdir = filepath.Join(p, gitdir)
				}
			}
			head, e := os.Stat(filepath.Join(gitdir, "HEAD"))
			if e != nil || !head.Mode().IsRegular() {
				return "", fmt.Errorf("invalid Git metadata: %s", marker)
			}
			common := gitdir
			if b, e := readBounded(filepath.Join(gitdir, "commondir"), 4096); e == nil {
				common = strings.TrimSpace(string(b))
				if !filepath.IsAbs(common) {
					common = filepath.Join(gitdir, common)
				}
			}
			objects, e := os.Stat(filepath.Join(common, "objects"))
			if e != nil || !objects.IsDir() {
				return "", fmt.Errorf("missing Git objects: %s", marker)
			}
			return p, nil
		}
		if !os.IsNotExist(e) {
			return "", e
		}
		if filepath.Dir(p) == p {
			break
		}
	}
	return "", nil
}

// Discover inspects selected and its immediate children only. Ancestor fallback
// is reserved for implicit CWD selection, and cannot broaden the downward scan.
func Discover(selected string, explicit bool) (Collection, error) {
	c := Collection{Projects: []Project{}, Diagnostics: []Diagnostic{}}
	abs, err := filepath.Abs(selected)
	if err != nil {
		return c, err
	}
	root, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return c, err
	}
	st, err := os.Stat(root)
	if err != nil {
		return c, err
	}
	if !st.IsDir() {
		return c, fmt.Errorf("selection must be a directory")
	}
	addDiag := func(code, p, msg string) {
		c.Diagnostics = append(c.Diagnostics, Diagnostic{Code: code, Severity: "warning", Path: p, Message: msg})
	}
	seen := map[string]bool{}
	inspect := func(p string, allowGeneric bool) bool {
		bundle := filepath.Join(p, "factory")
		st, e := os.Lstat(bundle)
		if os.IsNotExist(e) {
			return false
		}
		if e != nil {
			addDiag("inaccessible_path", bundle, e.Error())
			return false
		}
		if st.Mode()&os.ModeSymlink != 0 {
			addDiag("symlink_excluded", bundle, "bundle symlinks are not followed")
			return false
		}
		if !st.IsDir() {
			addDiag("invalid_bundle", bundle, "factory must be a directory")
			return false
		}
		repo, e := FindRepository(p)
		if e != nil {
			addDiag("invalid_repository", p, e.Error())
			return false
		}
		if repo == "" {
			addDiag("not_git_project", p, "bundle is outside a Git working tree")
			return false
		}
		rel, _ := filepath.Rel(repo, p)
		for _, component := range strings.Split(filepath.ToSlash(rel), "/") {
			if component == ".git" {
				addDiag("git_internals_excluded", p, "Git internals are not project roots")
				return false
			}
		}
		if seen[p] {
			return true
		}
		project, diags, hasFactory := metadata(p, repo)
		if !hasFactory && !allowGeneric {
			return false
		}
		seen[p] = true
		c.Projects = append(c.Projects, project)
		c.Diagnostics = append(c.Diagnostics, diags...)
		return true
	}
	inspect(root, explicit)
	entries, err := os.ReadDir(root)
	if err != nil {
		return c, err
	}
	repo, _ := FindRepository(root)
	for _, entry := range entries {
		p := filepath.Join(root, entry.Name())
		if entry.Name() == ".git" || entry.Name() == ".factory" || entry.Name() == "factory" {
			continue
		}
		if entry.Type()&os.ModeSymlink != 0 {
			addDiag("symlink_excluded", p, "directory symlinks are not followed")
			continue
		}
		if !entry.IsDir() {
			continue
		}
		if ignored(p, repo) {
			addDiag("ignored_path", p, "excluded by Git ignore patterns")
			continue
		}
		inspect(p, false)
	}
	if len(c.Projects) == 0 && !explicit && repo != "" {
		for p := filepath.Dir(root); within(repo, p); p = filepath.Dir(p) {
			if inspect(p, false) {
				break
			}
			if p == repo {
				break
			}
		}
	}
	ids := map[string]string{}
	uris := map[string]string{}
	for _, p := range c.Projects {
		if !strings.HasPrefix(p.ID, "local:") {
			if prior, ok := ids[p.ID]; ok {
				addDiag("duplicate_project_id", p.Root, "project identity also occurs at "+prior)
			} else {
				ids[p.ID] = p.Root
			}
		}
		if p.BundleURI != "" {
			if prior, ok := uris[p.BundleURI]; ok {
				addDiag("duplicate_bundle_uri", p.Root, "canonical bundle URL also occurs at "+prior)
			} else {
				uris[p.BundleURI] = p.Root
			}
		}
	}
	sort.Slice(c.Projects, func(i, j int) bool {
		if c.Projects[i].Name != c.Projects[j].Name {
			return c.Projects[i].Name < c.Projects[j].Name
		}
		return c.Projects[i].Root < c.Projects[j].Root
	})
	return c, nil
}

func metadata(root, repo string) (Project, []Diagnostic, bool) {
	h := fmt.Sprintf("%x", sha256.Sum256([]byte(root)))
	p := Project{ID: "local:" + h[:16], Checkout: h, Name: filepath.Base(root), Root: root, RepositoryRoot: repo}
	file := filepath.Join(root, "factory", "factory.md")
	bad := func(msg string) (Project, []Diagnostic, bool) {
		return p, []Diagnostic{{Code: "invalid_factory_metadata", Severity: "warning", Path: file, Message: msg}}, true
	}
	st, e := os.Lstat(file)
	if os.IsNotExist(e) {
		return p, nil, false
	}
	if e != nil {
		return bad(e.Error())
	}
	if !st.Mode().IsRegular() {
		return bad("Factory metadata must be a regular file, not a symlink")
	}
	b, e := readBounded(file, 1024*1024)
	if e != nil {
		return bad(e.Error())
	}
	if !utf8.Valid(b) {
		return bad("Factory metadata must be UTF-8")
	}
	s := strings.ReplaceAll(string(b), "\r\n", "\n")
	if !strings.HasPrefix(s, "---\n") {
		return bad("missing YAML frontmatter")
	}
	end := strings.Index(s[4:], "\n---\n")
	if end < 0 {
		return bad("missing frontmatter closing delimiter")
	}
	if end > 256*1024 {
		return bad("frontmatter exceeds 256 KiB")
	}
	var n yaml.Node
	dec := yaml.NewDecoder(strings.NewReader(s[4 : 4+end]))
	if e := dec.Decode(&n); e != nil {
		return bad("invalid YAML: " + e.Error())
	}
	if len(n.Content) != 1 || n.Content[0].Kind != yaml.MappingNode {
		return bad("frontmatter must be a mapping")
	}
	if e := safeYAML(n.Content[0], 0); e != nil {
		return bad(e.Error())
	}
	var fields map[string]any
	if e := n.Decode(&fields); e != nil {
		return bad(e.Error())
	}
	if fields["type"] != "Factory" {
		return bad("factory.md must have type Factory")
	}
	id, ok := fields["project_id"].(string)
	if !ok || strings.TrimSpace(id) == "" || strings.ContainsAny(id, "\r\n\t /\\") || strings.HasPrefix(id, "local:") {
		return bad("project_id must be a nonempty portable identifier")
	}
	name, ok := fields["title"].(string)
	if !ok || strings.TrimSpace(name) == "" {
		return bad("Factory title must be a nonempty string")
	}
	if raw, exists := fields["bundle_uri"]; exists {
		uri, ok := raw.(string)
		u, e := url.Parse(uri)
		if !ok || e != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || !strings.HasSuffix(uri, "/") {
			return bad("bundle_uri must be an absolute HTTPS URL prefix ending in /")
		}
		p.BundleURI = uri
	}
	p.ID = id
	p.Name = name
	return p, nil, true
}
func safeYAML(n *yaml.Node, depth int) error {
	if depth > 32 {
		return fmt.Errorf("YAML nesting exceeds 32")
	}
	if n.Kind == yaml.AliasNode || n.Anchor != "" {
		return fmt.Errorf("YAML aliases and anchors are unsupported")
	}
	if n.Tag != "" && !strings.HasPrefix(n.Tag, "!!") {
		return fmt.Errorf("custom YAML tags are unsupported")
	}
	if n.Kind == yaml.MappingNode {
		keys := map[string]bool{}
		for i := 0; i < len(n.Content); i += 2 {
			k := n.Content[i]
			if k.Kind != yaml.ScalarNode || k.Tag != "!!str" {
				return fmt.Errorf("YAML keys must be strings")
			}
			if keys[k.Value] {
				return fmt.Errorf("duplicate YAML key: %s", k.Value)
			}
			keys[k.Value] = true
		}
	}
	for _, v := range n.Content {
		if e := safeYAML(v, depth+1); e != nil {
			return e
		}
	}
	return nil
}
func readBounded(p string, max int64) ([]byte, error) {
	f, e := safeio.OpenRegular(p, os.O_RDONLY)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	b, e := io.ReadAll(io.LimitReader(f, max+1))
	if int64(len(b)) > max {
		return nil, fmt.Errorf("file exceeds size limit: %s", p)
	}
	return b, e
}
func within(root, p string) bool {
	rel, e := filepath.Rel(root, p)
	return e == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// Ignore rules are read from each enclosing .gitignore and .git/info/exclude.
// No global Git configuration or user preference files are consulted.
func ignored(p, repo string) bool {
	if repo == "" {
		return false
	}
	result := false
	apply := func(base, file string) {
		b, e := readBounded(file, 1024*1024)
		if e != nil {
			return
		}
		rel, e := filepath.Rel(base, p)
		if e != nil {
			return
		}
		rel = filepath.ToSlash(rel)
		for _, line := range strings.Split(string(b), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			neg := strings.HasPrefix(line, "!")
			if neg {
				line = line[1:]
			}
			line = strings.TrimSuffix(line, "/")
			anchored := strings.HasPrefix(line, "/")
			line = strings.TrimPrefix(line, "/")
			matched := false
			if anchored || strings.Contains(line, "/") {
				matched = glob(line, rel)
			} else {
				for _, part := range strings.Split(rel, "/") {
					if glob(line, part) {
						matched = true
						break
					}
				}
			}
			if matched {
				result = !neg
			}
		}
	}
	apply(repo, filepath.Join(repo, ".git", "info", "exclude"))
	dirs := []string{}
	for d := filepath.Dir(p); within(repo, d); d = filepath.Dir(d) {
		dirs = append(dirs, d)
		if d == repo {
			break
		}
	}
	for i := len(dirs) - 1; i >= 0; i-- {
		apply(dirs[i], filepath.Join(dirs[i], ".gitignore"))
	}
	return result
}
func glob(pattern, value string) bool {
	if strings.HasPrefix(pattern, "**/") {
		if glob(strings.TrimPrefix(pattern, "**/"), value) {
			return true
		}
		for i := strings.Index(value, "/"); i >= 0; i = strings.Index(value, "/") {
			value = value[i+1:]
			if glob(strings.TrimPrefix(pattern, "**/"), value) {
				return true
			}
		}
		return false
	}
	if strings.HasSuffix(pattern, "/**") {
		base := strings.TrimSuffix(pattern, "/**")
		return value == base || strings.HasPrefix(value, base+"/")
	}
	ok, _ := filepath.Match(pattern, value)
	return ok
}
