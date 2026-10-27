//go:build darwin || linux

package app_test

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/tclasen/factoryctl/internal/app"
	"golang.org/x/sys/unix"
)

// A separate OS process holds precisely the advisory lock used by the store.
// The parent waits for a pipe handshake, then verifies rejection and recovery.
func TestConcurrencyLockHelper(t *testing.T) {
	if os.Getenv("FACTORYCTL_LOCK_HELPER") != "1" {
		return
	}
	f, e := os.OpenFile(os.Getenv("FACTORYCTL_LOCK_ANCHOR"), os.O_RDWR, 0)
	if e != nil {
		t.Fatal(e)
	}
	defer f.Close()
	if e = unix.Flock(int(f.Fd()), unix.LOCK_EX); e != nil {
		t.Fatal(e)
	}
	fmt.Println("locked")
	var b [1]byte
	_, _ = os.Stdin.Read(b[:])
}
func TestConcurrencyProcessLockAndExit(t *testing.T) {
	projects := collectionFixture(t, "S1")
	p := projects[0]
	s, e := app.Load(projects)
	if e != nil {
		t.Fatal(e)
	}
	d := s.Bundle(p.Checkout).Docs["widgets/example.md"]
	file := filepath.Join(p.Root, "factory", d.Ref.Path)
	before := collectionRead(t, file)
	executable, e := os.Executable()
	if e != nil {
		t.Fatal(e)
	}
	cmd := exec.Command(executable, "-test.run=^TestConcurrencyLockHelper$", "-test.timeout=15s")
	cmd.Env = append(os.Environ(), "FACTORYCTL_LOCK_HELPER=1", "FACTORYCTL_LOCK_ANCHOR="+filepath.Join(p.Root, "factory/control/write-coordination.md"))
	stdout, e := cmd.StdoutPipe()
	if e != nil {
		t.Fatal(e)
	}
	stdin, e := cmd.StdinPipe()
	if e != nil {
		t.Fatal(e)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if e = cmd.Start(); e != nil {
		t.Fatal(e)
	}
	waited := false
	t.Cleanup(func() {
		_ = stdin.Close()
		if !waited {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		}
	})
	ready := make(chan string, 1)
	go func() { line, _ := bufio.NewReader(stdout).ReadString('\n'); ready <- line }()
	select {
	case line := <-ready:
		if line != "locked\n" {
			t.Fatalf("helper: %q %s", line, stderr.String())
		}
	case <-time.After(5 * time.Second):
		t.Fatal("lock helper failed to become ready")
	}
	_, e = app.Mutate(projects, p, app.Mutation{Kind: "update", Path: d.Ref.Path, ExpectedRevision: d.Revision, Set: map[string]any{"state": "alpha-ready"}})
	collectionCode(t, e, "stale_revision")
	if !bytes.Equal(before, collectionRead(t, file)) {
		t.Fatal("busy writer changed bytes")
	}
	// Abrupt process exit releases the kernel lock without any lock text repair.
	if e = cmd.Process.Kill(); e != nil {
		t.Fatal(e)
	}
	_ = cmd.Wait()
	waited = true
	changed, e := app.Mutate(projects, p, app.Mutation{Kind: "update", Path: d.Ref.Path, ExpectedRevision: d.Revision, Set: map[string]any{"state": "alpha-ready"}})
	if e != nil {
		t.Fatalf("lock not released on process exit: %v", e)
	}
	if changed.Revision == d.Revision {
		t.Fatal("successful write did not publish changed bytes")
	}
}
func TestConcurrencyExternalEditRevision(t *testing.T) {
	projects := collectionFixture(t, "S1")
	p := projects[0]
	s, e := app.Load(projects)
	if e != nil {
		t.Fatal(e)
	}
	d := s.Bundle(p.Checkout).Docs["widgets/example.md"]
	file := filepath.Join(p.Root, "factory", d.Ref.Path)
	external := append(append([]byte{}, d.Raw...), []byte("\nExternally edited.\n")...)
	collectionPut(t, file, external)
	_, e = app.Mutate(projects, p, app.Mutation{Kind: "replace", Path: d.Ref.Path, ExpectedRevision: d.Revision, Document: d.Raw})
	collectionCode(t, e, "stale_revision")
	if !bytes.Equal(external, collectionRead(t, file)) {
		t.Fatal("stale replacement overwrote external editor")
	}
}
func TestConcurrencyCreateExclusivity(t *testing.T) {
	projects := collectionFixture(t, "S1")
	p := projects[0]
	s, e := app.Load(projects)
	if e != nil {
		t.Fatal(e)
	}
	d := s.Bundle(p.Checkout).Docs["widgets/example.md"]
	created, e := app.Mutate(projects, p, app.Mutation{Kind: "create", Path: "widgets/new.md", Document: d.Raw})
	if e != nil {
		t.Fatal(e)
	}
	_, e = app.Mutate(projects, p, app.Mutation{Kind: "create", Path: "widgets/new.md", Document: append(append([]byte{}, d.Raw...), []byte("changed\n")...)})
	collectionCode(t, e, "already_exists")
	if !bytes.Equal(created.Raw, collectionRead(t, filepath.Join(p.Root, "factory/widgets/new.md"))) {
		t.Fatal("create replaced existing file")
	}
}
func TestConcurrencySymlinkWriteRejection(t *testing.T) {
	for _, kind := range []string{"bundle-root", "ancestor", "final"} {
		t.Run(kind, func(t *testing.T) {
			projects := collectionFixture(t, "S1")
			p := projects[0]
			s, e := app.Load(projects)
			if e != nil {
				t.Fatal(e)
			}
			d := s.Bundle(p.Checkout).Docs["widgets/example.md"]
			outside := t.TempDir()
			sentinel := filepath.Join(outside, "sentinel.md")
			collectionPut(t, sentinel, d.Raw)
			mutation := app.Mutation{Kind: "create", Path: "widgets/new.md", Document: d.Raw}
			switch kind {
			case "bundle-root":
				bundle := filepath.Join(p.Root, "factory")
				moved := filepath.Join(outside, "bundle")
				if e = os.Rename(bundle, moved); e != nil {
					t.Fatal(e)
				}
				if e = os.Symlink(moved, bundle); e != nil {
					t.Fatal(e)
				}
			case "ancestor":
				if e = os.Symlink(outside, filepath.Join(p.Root, "factory", "linked")); e != nil {
					t.Fatal(e)
				}
				mutation.Path = "linked/new.md"
			case "final":
				if e = os.Symlink(sentinel, filepath.Join(p.Root, "factory/widgets/symlink.md")); e != nil {
					t.Fatal(e)
				}
				mutation = app.Mutation{Kind: "replace", Path: "widgets/symlink.md", ExpectedRevision: d.Revision, Document: []byte("changed")}
			}
			_, e = app.Mutate(projects, p, mutation)
			collectionCode(t, e, "invalid_path")
			if !bytes.Equal(d.Raw, collectionRead(t, sentinel)) {
				t.Fatal("symlink mutation changed external file")
			}
			if _, e = os.Stat(filepath.Join(outside, "new.md")); !os.IsNotExist(e) {
				t.Fatalf("unexpected external create: %v", e)
			}
		})
	}
}
