package runtime

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func settingsStore(t *testing.T) (*ChatSettingsStore, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "personal", "chat-settings.json")
	s, e := NewChatSettingsStore(path)
	if e != nil {
		t.Fatal(e)
	}
	return s, path
}
func TestSettingsPersistenceResetIsolationAndPermissions(t *testing.T) {
	s, path := settingsStore(t)
	a, e := s.Read()
	if e != nil || a.Revision != "absent" {
		t.Fatal(a, e)
	}
	if _, e = os.Stat(filepath.Dir(path)); !os.IsNotExist(e) {
		t.Fatal("read wrote directory")
	}
	ctx := context.Background()
	a, e = s.Save(ctx, "", &ChatSelection{"global", "low"}, a.Revision)
	if e != nil {
		t.Fatal(e)
	}
	project := t.TempDir()
	other := t.TempDir()
	a, e = s.Save(ctx, project, &ChatSelection{"local", "high"}, a.Revision)
	if e != nil {
		t.Fatal(e)
	}
	a, e = s.Save(ctx, other, &ChatSelection{}, a.Revision)
	if e != nil {
		t.Fatal(e)
	}
	reloaded, _ := NewChatSettingsStore(path)
	b, e := reloaded.Read()
	if e != nil || b.Revision != a.Revision || len(b.Projects) != 2 {
		t.Fatal(b, e)
	}
	root, _ := CanonicalProjectRoot(project)
	if b.Projects[root].Model != "local" {
		t.Fatal(b)
	}
	b, e = s.Save(ctx, project, nil, b.Revision)
	if e != nil || len(b.Projects) != 1 || b.Global.Model != "global" {
		t.Fatal(b, e)
	}
	if _, e = s.Save(ctx, "", nil, a.Revision); !errors.Is(e, ErrSettingsConflict) {
		t.Fatal(e)
	}
	b, e = s.Save(ctx, "", nil, b.Revision)
	if e != nil || b.Global != nil {
		t.Fatal(b, e)
	}
	info, _ := os.Stat(path)
	if info.Mode().Perm() != 0600 {
		t.Fatal(info.Mode())
	}
	dir, _ := os.Stat(filepath.Dir(path))
	if dir.Mode().Perm() != 0700 {
		t.Fatal(dir.Mode())
	}
}
func TestSelectionPrecedence(t *testing.T) {
	global := ChatSelection{"global", "low"}
	saved := ChatSettingsSnapshot{Global: &global, Projects: map[string]ChatSelection{"/project": {"project", "high"}}}
	legacy := ChatSelection{"legacy", "medium"}
	tests := []struct {
		explicit *ChatSelection
		inv      ChatSelection
		root     string
		want     ChatSelection
	}{{nil, ChatSelection{}, "/project", ChatSelection{"project", "high"}}, {nil, ChatSelection{}, "/other", global}, {nil, ChatSelection{"chosen", ""}, "/project", ChatSelection{"chosen", ""}}, {nil, ChatSelection{"", "low"}, "/project", ChatSelection{"project", "low"}}, {&ChatSelection{}, ChatSelection{"chosen", "high"}, "/project", ChatSelection{}}}
	for _, tt := range tests {
		if got := ResolveChatSelection(tt.explicit, tt.inv, saved, tt.root, legacy); got != tt.want {
			t.Fatal(got, tt.want)
		}
	}
	if got := ResolveChatSelection(nil, ChatSelection{}, ChatSettingsSnapshot{}, "/x", legacy); got != legacy {
		t.Fatal(got)
	}
}
func TestSettingsStrictJSONAndSymlink(t *testing.T) {
	s, path := settingsStore(t)
	if e := os.Mkdir(filepath.Dir(path), 0700); e != nil {
		t.Fatal(e)
	}
	for _, bad := range []string{`null`, `{"version":1,"global":{"model":null,"reasoning_effort":"low"}}`, `{"version":1,"global":{"model":"a"}}`, `{"version":1,"projects":{"/project":null}}`, `{"version":1,"projects":{"/project":{"model":"a","reasoning_effort":null}}}`, `{"version":2}`, `{"version":1,"version":1}`, `{"version":1,"unknown":1}`, `{"version":1,"global":{"model":"a","model":"b"}}`, `{"version":1,"projects":{"relative":{"model":"a"}}}`, `{"version":1,"global":{"extra":"x"}}`, `{"version":1} {}`, strings.Repeat(" ", maxChatSettings+1)} {
		os.WriteFile(path, []byte(bad), 0600)
		if _, e := s.Read(); e == nil {
			t.Fatalf("accepted %s", bad)
		}
	}
	os.Remove(path)
	target := filepath.Join(t.TempDir(), "outside")
	os.WriteFile(target, []byte(`{"version":1}`), 0600)
	if e := os.Symlink(target, path); e != nil {
		t.Fatal(e)
	}
	if _, e := s.Read(); e == nil {
		t.Fatal("read symlink")
	}
	if _, e := s.Save(context.Background(), "", &ChatSelection{}, "absent"); e == nil {
		t.Fatal("saved symlink")
	}
	b, _ := os.ReadFile(target)
	if string(b) != `{"version":1}` {
		t.Fatal("changed target")
	}
}
func TestSettingsLockCancellation(t *testing.T) {
	s, path := settingsStore(t)
	a, e := s.Save(context.Background(), "", &ChatSelection{}, "absent")
	if e != nil {
		t.Fatal(e)
	}
	f, e := os.OpenFile(path+".lock", os.O_RDWR, 0600)
	if e != nil {
		t.Fatal(e)
	}
	defer f.Close()
	if e = lockChatSettings(f); e != nil {
		t.Fatal(e)
	}
	defer unlockChatSettings(f)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	if _, e = s.Save(ctx, "", nil, a.Revision); !errors.Is(e, context.DeadlineExceeded) {
		t.Fatal(e)
	}
}
func TestSettingsProcessHelper(t *testing.T) {
	if os.Getenv("FACTORYCTL_SETTINGS_HELPER") != "1" {
		return
	}
	s, e := NewChatSettingsStore(os.Getenv("FACTORYCTL_SETTINGS_PATH"))
	if e != nil {
		os.Exit(3)
	}
	_, e = s.Save(context.Background(), "", &ChatSelection{Model: os.Getenv("FACTORYCTL_SETTINGS_MODEL")}, "absent")
	if errors.Is(e, ErrSettingsConflict) {
		os.Exit(2)
	}
	if e != nil {
		os.Exit(3)
	}
	os.Exit(0)
}
func TestSettingsConcurrentProcesses(t *testing.T) {
	_, path := settingsStore(t)
	codes := make(chan int, 2)
	var wg sync.WaitGroup
	for _, name := range []string{"first", "second"} {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			cmd := exec.Command(os.Args[0], "-test.run=^TestSettingsProcessHelper$")
			cmd.Env = append(os.Environ(), "FACTORYCTL_SETTINGS_HELPER=1", "FACTORYCTL_SETTINGS_PATH="+path, "FACTORYCTL_SETTINGS_MODEL="+name)
			e := cmd.Run()
			if e == nil {
				codes <- 0
				return
			}
			var exit *exec.ExitError
			if errors.As(e, &exit) {
				codes <- exit.ExitCode()
			} else {
				codes <- 99
			}
		}(name)
	}
	wg.Wait()
	close(codes)
	success, conflict := 0, 0
	for n := range codes {
		if n == 0 {
			success++
		}
		if n == 2 {
			conflict++
		}
	}
	if success != 1 || conflict != 1 {
		t.Fatal(success, conflict)
	}
	b, e := os.ReadFile(path)
	if e != nil || !json.Valid(b) {
		t.Fatal(string(b), e)
	}
}

func TestSettingsRejectUnsafeDirectoryAndLock(t *testing.T) {
	s, path := settingsStore(t)
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Dir(path)); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Read(); err == nil {
		t.Fatal("accepted symlink settings directory")
	}
	if _, err := s.Save(context.Background(), "", &ChatSelection{}, "absent"); err == nil {
		t.Fatal("wrote through symlink directory")
	}
	os.Remove(filepath.Dir(path))
	os.Mkdir(filepath.Dir(path), 0700)
	target := filepath.Join(outside, "lock-target")
	os.WriteFile(target, []byte("unchanged"), 0600)
	if err := os.Symlink(target, path+".lock"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Save(context.Background(), "", &ChatSelection{}, "absent"); err == nil {
		t.Fatal("accepted symlink lock")
	}
	bytes, _ := os.ReadFile(target)
	if string(bytes) != "unchanged" {
		t.Fatal("changed external target")
	}
}
func TestSettingsReadOnlyDirectoryError(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root bypasses directory permissions")
	}
	s, path := settingsStore(t)
	if err := os.Mkdir(filepath.Dir(path), 0500); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(filepath.Dir(path), 0700)
	if _, err := s.Save(context.Background(), "", &ChatSelection{}, "absent"); err == nil {
		t.Fatal("read-only settings directory was writable")
	}
}
