package runtime

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/tclasen/factoryctl/internal/safeio"
	"gopkg.in/yaml.v3"
)

// Selection is contributor-owned launch data. Call ResolvePreferences only for
// an explicit runtime action; ordinary document access never reads preferences.
type Selection struct {
	Runtime     string `json:"runtime"`
	Executable  string `json:"executable,omitempty"`
	Model       string `json:"model,omitempty"`
	Effort      string `json:"reasoning_effort,omitempty"`
	Preferences string `json:"-"`
}

func ResolvePreferences(explicit Selection) (Selection, error) {
	selectedFile := explicit.Preferences
	if selectedFile == "" {
		u, err := user.Current()
		if err != nil {
			return Selection{}, errors.New("cannot resolve current user's home for runtime preferences")
		}
		selectedFile = filepath.Join(u.HomeDir, ".factoryctl", "preferences.md")
	}
	b, err := safeio.ReadFile(selectedFile, 1<<20)
	var file Selection
	if err != nil {
		if !(explicit.Preferences == "" && errors.Is(err, os.ErrNotExist)) {
			return Selection{}, fmt.Errorf("runtime preferences unavailable: %w", err)
		}
	} else {
		if len(b) > 1<<20 {
			return Selection{}, errors.New("runtime preferences exceed size limit")
		}
		file, err = parsePreferences(b)
		if err != nil {
			return Selection{}, err
		}
	}
	out := file
	if explicit.Runtime != "" && explicit.Runtime != file.Runtime {
		out = Selection{Runtime: explicit.Runtime}
	}
	if explicit.Runtime != "" {
		out.Runtime = explicit.Runtime
	}
	if explicit.Executable != "" {
		out.Executable = explicit.Executable
	}
	if explicit.Model != "" {
		out.Model = explicit.Model
	}
	if explicit.Effort != "" {
		out.Effort = explicit.Effort
	}
	out.Preferences = selectedFile
	if out.Runtime == "" {
		return Selection{}, errors.New("select a runtime explicitly or in personal preferences")
	}
	if out.Runtime != "codex" {
		return Selection{}, errors.New("unsupported runtime adapter; no fallback selected")
	}
	if out.Executable == "" {
		out.Executable = "codex"
	} else if !filepath.IsAbs(out.Executable) {
		return Selection{}, errors.New("explicit runtime executable must be an absolute path")
	}
	return out, nil
}

func parsePreferences(b []byte) (Selection, error) {
	var s Selection
	if !utf8.Valid(b) {
		return s, errors.New("runtime preferences are not UTF-8")
	}
	b = bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n"))
	if !bytes.HasPrefix(b, []byte("---\n")) {
		return s, errors.New("runtime preferences require YAML frontmatter")
	}
	end := bytes.Index(b[4:], []byte("\n---"))
	if end < 0 {
		return s, errors.New("runtime preferences frontmatter is not closed")
	}
	end += 4
	if end > 256<<10 {
		return s, errors.New("runtime preferences frontmatter exceeds size limit")
	}
	if len(b) > end+4 && b[end+4] != '\n' {
		return s, errors.New("invalid runtime preferences frontmatter delimiter")
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(b[4:end], &doc); err != nil {
		return s, errors.New("invalid runtime preference YAML")
	}
	if len(doc.Content) != 1 || doc.Content[0].Kind != yaml.MappingNode || doc.Content[0].Tag != "!!map" || doc.Content[0].Anchor != "" {
		return s, errors.New("runtime preferences must be a mapping")
	}
	seen := map[string]bool{}
	version := false
	for i := 0; i < len(doc.Content[0].Content); i += 2 {
		k, v := doc.Content[0].Content[i], doc.Content[0].Content[i+1]
		if k.Kind != yaml.ScalarNode || k.Tag != "!!str" || k.Anchor != "" || v.Anchor != "" {
			return s, errors.New("invalid runtime preference key or anchor")
		}
		if seen[k.Value] {
			return s, errors.New("duplicate runtime preference field")
		}
		seen[k.Value] = true
		if k.Value == "version" {
			if v.Kind != yaml.ScalarNode || v.Tag != "!!int" || v.Value != "1" {
				return s, errors.New("runtime preferences require version 1")
			}
			version = true
			continue
		}
		if v.Kind != yaml.ScalarNode || v.Tag != "!!str" || strings.TrimSpace(v.Value) == "" {
			return s, errors.New("runtime preference values must be nonempty strings")
		}
		switch k.Value {
		case "runtime":
			s.Runtime = v.Value
		case "executable":
			s.Executable = v.Value
		case "model":
			s.Model = v.Value
		case "reasoning_effort":
			s.Effort = v.Value
		default:
			return s, errors.New("unknown runtime preference field")
		}
	}
	if !version {
		return s, errors.New("runtime preferences require version 1")
	}
	if s.Runtime == "" && (s.Executable != "" || s.Model != "" || s.Effort != "") {
		return s, errors.New("adapter-bound preferences require runtime")
	}
	return s, nil
}
