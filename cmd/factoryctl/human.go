package main

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/tclasen/factoryctl/internal/app"
	"github.com/tclasen/factoryctl/internal/knowledge"
	"github.com/tclasen/factoryctl/internal/project"
)

func human(out, errOut io.Writer, value any) {
	diagnostics := func(ds []knowledge.Diagnostic) {
		for _, d := range ds {
			fmt.Fprintf(errOut, "%s/%s:%d: %s: %s\n", d.ProjectID, d.Path, d.Line, d.Code, d.Message)
		}
	}
	nodes := func(ds []*knowledge.Document) {
		w := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
		fmt.Fprintln(w, "PROJECT\tTYPE\tPATH\tTITLE")
		for _, d := range ds {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", d.Ref.ProjectID, d.Type, d.Ref.Path, d.Title)
		}
		_ = w.Flush()
	}
	switch v := value.(type) {
	case string:
		fmt.Fprintln(out, v)
	case project.Collection:
		w := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tPROJECT ROOT\tREPOSITORY ROOT")
		for _, p := range v.Projects {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.ID, p.Name, p.Root, p.RepositoryRoot)
		}
		_ = w.Flush()
		for _, d := range v.Diagnostics {
			fmt.Fprintf(errOut, "%s: %s: %s\n", d.Path, d.Code, d.Message)
		}
		fmt.Fprintf(out, "%d project(s)\n", len(v.Projects))
	case app.Validation:
		fmt.Fprintf(out, "%s validation: %d errors, %d warnings, %d unchecked links across %d project(s)\n", v.Level, v.Errors, v.Warnings, v.UncheckedLinks, len(v.Scope))
		diagnostics(v.Diagnostics)
	case app.InitResult:
		fmt.Fprintf(out, "Initialized %s: %d created, %d unchanged\n", v.ProjectRoot, len(v.Created), len(v.Unchanged))
		for _, p := range v.Created {
			fmt.Fprintln(out, "  "+p)
		}
	case *knowledge.Document:
		fmt.Fprintf(out, "%s: %s (%s)\nRevision: %s\n", v.Ref.Path, v.Title, v.Type, v.Revision)
		diagnostics(v.Diagnostics)
	case []*knowledge.Document:
		nodes(v)
	case []map[string]string:
		for _, row := range v {
			fmt.Fprintf(out, "%-10s %s\n", row["status"], row["path"])
		}
	case map[string]any:
		if ds, ok := v["nodes"].([]*knowledge.Document); ok {
			nodes(ds)
			if cursor, _ := v["next_cursor"].(string); cursor != "" {
				fmt.Fprintln(out, "Next cursor: "+cursor)
			}
			if ds, ok := v["diagnostics"].([]knowledge.Diagnostic); ok {
				diagnostics(ds)
			}
			return
		}
		if d, ok := v["node"].(*knowledge.Document); ok {
			fmt.Fprintf(out, "%s — %s (%s)\nRevision: %s\n\n%s", d.Ref.Path, d.Title, d.Type, d.Revision, d.Raw)
			if links, ok := v["backlinks"].([]app.Link); ok {
				fmt.Fprintf(out, "\n%d loaded backlink(s)\n", len(links))
			}
			return
		}
		keys := []string{}
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(out, "%s: %v\n", k, v[k])
		}
	default:
		fmt.Fprintln(out, v)
	}
}
