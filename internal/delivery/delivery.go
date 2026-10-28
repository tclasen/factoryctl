// Package delivery presents attributable assessments and content freshness.
// It never evaluates project policy, invokes Git, or infers authority.
package delivery

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/tclasen/factoryctl/internal/knowledge"
	"github.com/tclasen/factoryctl/internal/safeio"
)

type Target struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}
type Observation struct {
	Version            int               `json:"version"`
	Record             string            `json:"record"`
	Target             Target            `json:"target"`
	CriteriaRevision   string            `json:"criteria_revision"`
	AssessmentRevision string            `json:"assessment_revision"`
	Files              map[string]string `json:"files"`
	ObservedAt         string            `json:"observed_at"`
	Observer           string            `json:"observer"`
	Evidence           []string          `json:"evidence"`
}
type Status struct {
	State         string   `json:"state"`
	Explanation   string   `json:"explanation"`
	Record        string   `json:"record"`
	Target        Target   `json:"target"`
	Assessor      string   `json:"assessor,omitempty"`
	AcceptedClaim bool     `json:"accepted_claim"`
	ObservedAt    string   `json:"observed_at,omitempty"`
	Observer      string   `json:"observer,omitempty"`
	Evidence      []string `json:"evidence,omitempty"`
	Criteria      string   `json:"criteria,omitempty"`
	Assessment    string   `json:"assessment,omitempty"`
}

// read permits regular contained files only, including when a parent is swapped
// by another editor. Root-relative opens provide an additional containment bound.
func read(root, file string) ([]byte, error) {
	if file == "" || path.IsAbs(file) || path.Clean(file) != file || file == ".." || strings.HasPrefix(file, "../") || strings.ContainsAny(file, "\\:\x00") {
		return nil, fmt.Errorf("invalid contained file path")
	}
	current := root
	for _, part := range append([]string{""}, strings.Split(file, "/")...) {
		current = filepath.Join(current, part)
		st, err := os.Lstat(current)
		if err != nil {
			return nil, err
		}
		if st.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("symlink is not evidence")
		}
	}
	r, err := os.OpenRoot(root)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	b, err := safeio.ReadRoot(r, file, knowledge.MaxDocument)
	if len(b) > knowledge.MaxDocument {
		return nil, fmt.Errorf("evidence exceeds document limit")
	}
	return b, err
}

// Resolve uses project-owned Markdown and optional ignored verification evidence.
// delivery_ref is a bundle-relative reference to the same assessment-owning task.
// Observation JSON is disposable evidence, never configuration or durable state.
func Resolve(root, record string) Status {
	s := Status{State: "Delivery criteria not defined", Record: record, Explanation: "No project-owned delivery criteria are recorded."}
	load := func(p string) (*knowledge.Document, error) {
		raw, err := read(root, "factory/"+p)
		if err != nil {
			return nil, err
		}
		d := knowledge.Parse(p, raw)
		if len(d.Diagnostics) > 0 {
			return nil, fmt.Errorf("invalid evidence Markdown")
		}
		return d, nil
	}
	d, err := load(record)
	if err != nil {
		s.State = "Delivery unverified"
		s.Explanation = err.Error()
		return s
	}
	originalPath, originalRevision := record, d.Revision
	if ref := knowledge.String(d.Frontmatter, "delivery_ref"); ref != "" {
		if path.Clean(ref) != ref || strings.HasPrefix(ref, "../") || path.IsAbs(ref) {
			s.State = "Delivery unverified"
			s.Explanation = "Invalid delivery_ref."
			return s
		}
		record = ref
		s.Record = ref
		d, err = load(ref)
		if err != nil {
			s.State = "Delivery unverified"
			s.Explanation = "Referenced delivery record is unavailable."
			return s
		}
		if knowledge.String(d.Frontmatter, "delivery_ref") != "" {
			s.State = "Delivery unverified"
			s.Explanation = "Chained delivery references are unsupported."
			return s
		}
	}
	m := knowledge.Mapping(d.Frontmatter["delivery"])
	if m == nil {
		return s
	}
	s.Criteria = knowledge.String(m, "criteria")
	s.Assessment = knowledge.String(m, "assessment")
	t := knowledge.Mapping(m["target"])
	s.Target = Target{knowledge.String(t, "kind"), knowledge.String(t, "name")}
	if s.Criteria == "" || s.Assessment == "" || s.Target.Name == "" {
		s.State = "Delivery unverified"
		s.Explanation = "Delivery mapping is incomplete: criteria, assessment and target are required."
		return s
	}
	s.State = "Delivery unverified"
	if s.Target.Kind != "local-working-tree" && s.Target.Kind != "local-branch" && s.Target.Kind != "remote-branch" {
		s.Explanation = "Unsupported explicit delivery target."
		return s
	}
	criteria, err := load(s.Criteria)
	if err != nil {
		s.Explanation = "Delivery criteria are unavailable."
		return s
	}
	assessment, err := load(s.Assessment)
	if err != nil {
		s.Explanation = "Acceptance assessment is unavailable."
		return s
	}
	s.Assessor = knowledge.String(assessment.Frontmatter, "assessor")
	s.AcceptedClaim, _ = assessment.Frontmatter["accepted"].(bool)
	// Stable key contains no path input and never overlaps verified project content.
	store := ".factory/delivery/" + knowledge.Hash([]byte(record)) + ".json"
	raw, err := read(root, store)
	if err != nil {
		s.Explanation = "No readable local observation; acceptance remains a recorded claim."
		return s
	}
	var o Observation
	if json.Unmarshal(raw, &o) != nil || o.Version != 1 || o.Record != record || o.Target != s.Target || o.Observer == "" || len(o.Evidence) == 0 || len(o.Files) == 0 {
		s.Explanation = "Observation identity, manifest or provenance is invalid."
		return s
	}
	s.ObservedAt = o.ObservedAt
	s.Observer = o.Observer
	s.Evidence = o.Evidence
	timestamp, err := time.Parse(time.RFC3339, o.ObservedAt)
	if err != nil || timestamp.After(time.Now().Add(time.Minute)) {
		s.Explanation = "Invalid observation timestamp."
		return s
	}
	if o.CriteriaRevision != criteria.Revision || o.AssessmentRevision != assessment.Revision || knowledge.String(assessment.Frontmatter, "criteria_revision") != criteria.Revision || s.Assessor == "" {
		s.State = "Changes since verification"
		s.Explanation = "Criteria or recorded assessment no longer matches this observation."
		return s
	}
	if o.Files["factory/"+s.Criteria] != criteria.Revision || o.Files["factory/"+s.Assessment] != assessment.Revision {
		s.Explanation = "Manifest must bind criteria and assessment content."
		return s
	}
	if s.Target.Kind != "local-working-tree" {
		s.Explanation = "External target observation recorded; branch content and authority are not independently verified by factoryctl."
		return s
	}
	check := func() (map[string]string, error) {
		actual := map[string]string{}
		for file, expected := range o.Files {
			if strings.HasPrefix(file, ".factory/delivery/") || strings.HasPrefix(file, ".git/") || file == ".git" {
				return nil, fmt.Errorf("observation store and Git metadata cannot verify themselves")
			}
			b, e := read(root, file)
			if expected == "deleted" && os.IsNotExist(e) {
				// A missing parent is acceptable only if its existing ancestors are safe.
				for parent := path.Dir(file); parent != "."; parent = path.Dir(parent) {
					if st, se := os.Lstat(filepath.Join(root, parent)); se == nil && st.Mode()&os.ModeSymlink != 0 {
						return nil, fmt.Errorf("symlink deletion ancestor")
					}
				}
				actual[file] = "deleted"
				continue
			}
			if e != nil {
				return nil, e
			}
			actual[file] = knowledge.Hash(b)
			if actual[file] != expected {
				return nil, fmt.Errorf("relevant content changed: %s", file)
			}
		}
		return actual, nil
	}
	first, err := check()
	if err != nil {
		s.State = "Changes since verification"
		s.Explanation = err.Error()
		return s
	}
	second, err := check()
	if err != nil || !reflect.DeepEqual(first, second) {
		s.Explanation = "Content changed during verification; refresh before relying on it."
		return s
	}
	original, originalErr := load(originalPath)
	owner, ownerErr := load(record)
	if originalErr != nil || ownerErr != nil || original.Revision != originalRevision || owner.Revision != d.Revision {
		s.State = "Changes since verification"
		s.Explanation = "Delivery record or plan reference changed during verification."
		return s
	}
	if s.AcceptedClaim {
		s.State = "Assessed delivered locally — content verified"
		s.Explanation = "Recorded acceptance claim with matching scoped content; assessor authority is not inferred."
	} else {
		s.State = "Acceptance not recorded"
		s.Explanation = "Manifest matches; the assessment does not record acceptance."
	}
	return s
}
