package main

import (
	"strings"
	"testing"
)

func TestWebLauncherKeepsScopeAndRejectsUnsafeModes(t *testing.T) {
	root := initialized(t)
	r := invoke(t, "", "--path", root, "web", "--listen", "0.0.0.0:8080")
	if r.code == 0 || !strings.Contains(r.stderr, "loopback") {
		t.Fatalf("external bind was not rejected: %+v", r)
	}
	r = invoke(t, "", "--path", root, "web", "--json")
	wantError(t, r, 2, "invalid_arguments")
	r = invoke(t, "", "--path", root, "web", "extra", "--json")
	wantError(t, r, 2, "invalid_arguments")
	r = invoke(t, "", "--path", root, "doctor", "--listen", "127.0.0.1:8080", "--json")
	wantError(t, r, 2, "invalid_arguments")
}
