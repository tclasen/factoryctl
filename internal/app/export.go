package app

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/tclasen/factoryctl/internal/defaults"
	"github.com/tclasen/factoryctl/internal/project"
)

// ExportTemplates materializes only absent defaults. Existing customizations are
// left byte-identical; initialization and upgrades never replace them silently.
func ExportTemplates(projects []project.Project, p project.Project, kind string) (map[string][]string, error) {
	result := map[string][]string{"created": {}, "unchanged": {}}
	if kind != "record" && kind != "prompt" && kind != "view" {
		return result, Fail("invalid_arguments", "kind must be record, prompt or view")
	}
	if kind == "view" {
		name := ".factory/views/generic.html"
		// Reject symlink ancestors even when they point inside the project.
		for _, dir := range []string{".factory", ".factory/views"} {
			full := filepath.Join(p.Root, dir)
			if st, e := os.Lstat(full); e == nil {
				if !st.IsDir() || st.Mode()&os.ModeSymlink != 0 {
					return result, Fail("invalid_path", "asset directory is not a regular directory")
				}
			} else if !os.IsNotExist(e) {
				return result, e
			}
		}
		root, e := os.OpenRoot(p.Root)
		if e != nil {
			return result, e
		}
		defer root.Close()
		if e = root.MkdirAll(".factory/views", 0755); e != nil {
			return result, e
		}
		f, e := root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if os.IsExist(e) {
			st, se := root.Lstat(name)
			if se != nil || !st.Mode().IsRegular() {
				return result, Fail("invalid_path", "existing asset is not a regular file")
			}
			result["unchanged"] = append(result["unchanged"], name)
			return result, nil
		}
		if e != nil {
			return result, e
		}
		_, e = f.WriteString("<!doctype html><html lang=\"en\"><meta charset=\"utf-8\"><title>Custom record view</title><body><h1>Repository view</h1><p>Enable this view explicitly. Records arrive through the scoped message bridge.</p></body></html>\n")
		if e == nil {
			e = f.Sync()
		}
		ce := f.Close()
		if e == nil {
			e = ce
		}
		if e != nil {
			return result, e
		}
		result["created"] = append(result["created"], name)
		return result, nil
	}
	files := defaults.Files(true, p.ID, p.Name)
	keys := []string{}
	prefix := "templates/"
	if kind == "prompt" {
		prefix = "prompts/"
	}
	for file := range files {
		if strings.HasPrefix(file, prefix) {
			keys = append(keys, file)
		}
	}
	sort.Strings(keys)
	for _, file := range keys {
		if st, e := os.Lstat(filepath.Join(p.Root, "factory", file)); e == nil {
			if !st.Mode().IsRegular() {
				return result, Fail("invalid_path", "existing template is not a regular file")
			}
			result["unchanged"] = append(result["unchanged"], file)
			continue
		} else if !os.IsNotExist(e) {
			return result, e
		}
		if _, e := Mutate(projects, p, Mutation{Kind: "create", Path: file, Document: files[file]}); e != nil {
			return result, e
		}
		result["created"] = append(result["created"], file)
	}
	return result, nil
}
