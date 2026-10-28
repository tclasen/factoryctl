package runtime

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPreferenceSelectionAndAdapterIsolation(t *testing.T) {
	p := filepath.Join(t.TempDir(), "preferences.md")
	content := []byte("---\nversion: 1\nruntime: other\nexecutable: /other/runtime\nmodel: other-model\nreasoning_effort: other-effort\n---\nPersonal settings.\n")
	if err := os.WriteFile(p, content, 0600); err != nil {
		t.Fatal(err)
	}
	s, err := ResolvePreferences(Selection{Runtime: "codex", Preferences: p})
	if err != nil {
		t.Fatal(err)
	}
	if s.Executable != "codex" || s.Model != "" || s.Effort != "" {
		t.Fatal("adapter-bound fields leaked", s)
	}
	if _, err = ResolvePreferences(Selection{Preferences: p}); err == nil {
		t.Fatal("unsupported adapter silently substituted")
	}
	content = []byte("---\nversion: 1\nruntime: codex\nmodel: selected\nreasoning_effort: low\n---\n")
	os.WriteFile(p, content, 0600)
	s, err = ResolvePreferences(Selection{Preferences: p, Model: "explicit"})
	if err != nil {
		t.Fatal(err)
	}
	if s.Model != "explicit" || s.Effort != "low" {
		t.Fatal(s)
	}
	got, _ := os.ReadFile(p)
	if string(got) != string(content) {
		t.Fatal("preferences were modified")
	}
	if _, err = ResolvePreferences(Selection{Runtime: "codex", Preferences: filepath.Join(t.TempDir(), "absent")}); err == nil {
		t.Fatal("missing explicit preferences accepted")
	}
}
func TestMalformedPreferencesFailWithoutExecution(t *testing.T) {
	cases := []string{"version: 1\nruntime: codex", "---\nversion: 2\n---\n", "---\nversion: 1\nruntime: codex\nruntime: codex\n---\n", "---\nversion: 1\nruntime: &a codex\nmodel: *a\n---\n", "---\nversion: 1\nruntime: !custom codex\n---\n", "---\nversion: 1\nmodel: x\n---\n", "---\nversion: 1\nunknown: x\n---\n", "---\nversion: 1\nruntime: ''\n---\n"}
	for _, input := range cases {
		if _, err := parsePreferences([]byte(input)); err == nil {
			t.Fatalf("accepted malformed preferences %q", input)
		}
	}
}
