//go:build darwin || linux

package main

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestFIFOCommandHelper(t *testing.T) {
	raw := os.Getenv("FACTORYCTL_FIFO_ARGS")
	if raw == "" {
		return
	}
	var args []string
	if e := json.Unmarshal([]byte(raw), &args); e != nil {
		os.Exit(90)
	}
	os.Exit(run(args, strings.NewReader(""), os.Stdout, os.Stderr))
}
func TestSpecialFileReadsTerminate(t *testing.T) {
	for _, kind := range []string{"document", "asset", "ignore", "metadata"} {
		t.Run(kind, func(t *testing.T) {
			root := t.TempDir()
			for _, dir := range []string{".git/objects", ".git/refs", "factory/templates", ".factory/templates"} {
				if e := os.MkdirAll(filepath.Join(root, dir), 0755); e != nil {
					t.Fatal(e)
				}
			}
			if e := os.WriteFile(filepath.Join(root, ".git/HEAD"), []byte("ref: refs/heads/main\n"), 0644); e != nil {
				t.Fatal(e)
			}
			if e := os.WriteFile(filepath.Join(root, "factory/example.md"), []byte("---\ntype: Custom\n---\n"), 0644); e != nil {
				t.Fatal(e)
			}
			fifo := "factory/pipe.md"
			switch kind {
			case "asset":
				fifo = ".factory/templates/pipe.md.tmpl"
				body := "---\ntype: Record Template\nfor_type: Custom\ninputs: {}\nasset: .factory/templates/pipe.md.tmpl\n---\n# Template\n"
				if e := os.WriteFile(filepath.Join(root, "factory/templates/custom.md"), []byte(body), 0644); e != nil {
					t.Fatal(e)
				}
			case "ignore":
				fifo = ".gitignore"
			case "metadata":
				fifo = "factory/factory.md"
			}
			if e := syscall.Mkfifo(filepath.Join(root, fifo), 0600); e != nil {
				t.Fatal(e)
			}
			args, _ := json.Marshal([]string{"--path", root, "validate", "--json"})
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestFIFOCommandHelper$")
			cmd.Env = append(os.Environ(), "FACTORYCTL_FIFO_ARGS="+string(args))
			output, e := cmd.CombinedOutput()
			if ctx.Err() != nil {
				t.Fatalf("%s read blocked beyond timeout: %s", kind, output)
			}
			if e != nil {
				if exit, ok := e.(*exec.ExitError); !ok || exit.ExitCode() != 1 {
					t.Fatalf("unexpected process result: %v %s", e, output)
				}
			}
			var result map[string]any
			if e := json.Unmarshal(output, &result); e != nil {
				t.Fatalf("not JSON: %s: %v", output, e)
			}
			if (kind == "document" || kind == "asset" || kind == "metadata") && result["errors"] == float64(0) {
				t.Fatalf("special-file error missing: %s", output)
			}
		})
	}
}
