package delivery

import (
	"encoding/json"
	"github.com/tclasen/factoryctl/internal/knowledge"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLocalDeliveryFreshnessAndSharedPlan(t *testing.T) {
	root := t.TempDir()
	put := func(p string, b []byte) {
		t.Helper()
		file := filepath.Join(root, p)
		if e := os.MkdirAll(filepath.Dir(file), 0755); e != nil {
			t.Fatal(e)
		}
		if e := os.WriteFile(file, b, 0644); e != nil {
			t.Fatal(e)
		}
	}
	markdown := func(m map[string]any) []byte {
		t.Helper()
		b, e := knowledge.Encode(m, "\n# Evidence\n")
		if e != nil {
			t.Fatal(e)
		}
		return b
	}
	criteria := markdown(map[string]any{"type": "Reference", "title": "Criteria"})
	assessment := markdown(map[string]any{"type": "Reference", "title": "Assessment", "assessor": "human:fixture", "accepted": true, "criteria_revision": knowledge.Hash(criteria)})
	record := "tasks/item.md"
	task := markdown(map[string]any{"type": "Task", "title": "Item", "delivery": map[string]any{"criteria": "sources/criteria.md", "assessment": "sources/assessment.md", "target": map[string]any{"kind": "local-working-tree", "name": "my checkout"}}})
	put("factory/"+record, task)
	put("factory/sources/criteria.md", criteria)
	put("factory/sources/assessment.md", assessment)
	put("product.txt", []byte("candidate"))
	put("factory/designs/plan.md", markdown(map[string]any{"type": "Design", "title": "Plan", "delivery_ref": record}))
	if s := Resolve(root, record); s.State != "Delivery unverified" {
		t.Fatal(s)
	}
	o := Observation{Version: 1, Record: record, Target: Target{"local-working-tree", "my checkout"}, CriteriaRevision: knowledge.Hash(criteria), AssessmentRevision: knowledge.Hash(assessment), Files: map[string]string{"factory/sources/criteria.md": knowledge.Hash(criteria), "factory/sources/assessment.md": knowledge.Hash(assessment), "product.txt": knowledge.Hash([]byte("candidate")), "removed.txt": "deleted"}, ObservedAt: time.Now().UTC().Format(time.RFC3339), Observer: "fixture:content-check", Evidence: []string{"local fixture acceptance"}}
	publish := func() { b, _ := json.Marshal(o); put(".factory/delivery/"+knowledge.Hash([]byte(record))+".json", b) }
	publish()
	a, b := Resolve(root, record), Resolve(root, "designs/plan.md")
	{
		aj, _ := json.Marshal(a)
		bj, _ := json.Marshal(b)
		if !strings.Contains(a.State, "content verified") || string(aj) != string(bj) {
			t.Fatalf("task=%+v plan=%+v", a, b)
		}
	}
	put("product.txt", []byte("new draft"))
	if s := Resolve(root, record); s.State != "Changes since verification" {
		t.Fatal(s)
	}
	put("product.txt", []byte("candidate"))
	put("removed.txt", []byte("unexpected"))
	if s := Resolve(root, record); s.State != "Changes since verification" {
		t.Fatal(s)
	}
	os.Remove(filepath.Join(root, "removed.txt"))
	o.Files[".factory/delivery/unsafe.json"] = "deleted"
	publish()
	if s := Resolve(root, record); strings.Contains(s.State, "content verified") {
		t.Fatal(s)
	}
}

func TestDeliveryDoesNotInferApprovalOrRemoteEvidence(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "factory/tasks"), 0755)
	os.WriteFile(filepath.Join(root, "factory/tasks/done.md"), []byte("---\ntype: Task\ntitle: Done\ntask_state: done\n---\n"), 0644)
	if s := Resolve(root, "tasks/done.md"); s.State != "Delivery criteria not defined" || s.AcceptedClaim {
		t.Fatal(s)
	}
	if s := Resolve(root, "../outside.md"); strings.Contains(s.State, "verified —") {
		t.Fatal(s)
	}
	outside := t.TempDir()
	os.WriteFile(filepath.Join(outside, "secret.md"), []byte("secret"), 0644)
	os.Symlink(outside, filepath.Join(root, "factory/linked"))
	if _, e := read(root, "factory/linked/secret.md"); e == nil {
		t.Fatal("followed symlink")
	}
}
