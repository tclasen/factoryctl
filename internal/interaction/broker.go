// Package interaction owns ephemeral, scope-bound questions and generated views.
// It does not interpret an answer as workflow approval or persist conversations.
package interaction

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/tclasen/factoryctl/internal/app"
)

const MaxHTML = 128 << 10
const MaxText = 16 << 10
const MaxQuestions = 128
const MaxViews = 16
const MaxScopes = 32
const QuestionLifetime = 30 * time.Minute

type Choice struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}
type AskInput struct {
	Prompt            string   `json:"prompt"`
	Kind              string   `json:"kind"`
	Choices           []Choice `json:"choices,omitempty"`
	PreviousRequestID string   `json:"previous_request_id,omitempty"`
}
type AnswerInput struct {
	Status    string   `json:"status"`
	Text      string   `json:"text,omitempty"`
	OptionIDs []string `json:"option_ids,omitempty"`
}
type Question struct {
	ID                string    `json:"id"`
	ScopeID           string    `json:"scope_id"`
	Prompt            string    `json:"prompt"`
	Kind              string    `json:"kind"`
	Choices           []Choice  `json:"choices,omitempty"`
	PreviousRequestID string    `json:"previous_request_id,omitempty"`
	Status            string    `json:"status"`
	Text              string    `json:"text,omitempty"`
	OptionIDs         []string  `json:"option_ids,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	ExpiresAt         time.Time `json:"expires_at"`
}
type ViewInput struct {
	Title     string          `json:"title"`
	Fallback  string          `json:"fallback"`
	HTML      string          `json:"html"`
	Data      json.RawMessage `json:"data,omitempty"`
	RequestID string          `json:"request_id,omitempty"`
}
type View struct {
	ID          string          `json:"id"`
	ResourceURI string          `json:"resource_uri"`
	Title       string          `json:"title"`
	Fallback    string          `json:"fallback"`
	HTML        string          `json:"html"`
	Data        json.RawMessage `json:"data,omitempty"`
	RequestID   string          `json:"request_id,omitempty"`
}
type Snapshot struct {
	ScopeID   string     `json:"scope_id"`
	ProjectID string     `json:"project_id"`
	Closed    bool       `json:"closed"`
	Questions []Question `json:"questions"`
	Views     []View     `json:"views"`
}
type entry struct {
	q    Question
	done chan struct{}
}
type scope struct {
	project, token string
	closed         bool
	questions      []*entry
	views          []View
}
type Broker struct {
	mu     sync.Mutex
	scopes map[string]*scope
	tokens map[string]string
	closed bool
	now    func() time.Time
}

func New() *Broker {
	return &Broker{scopes: map[string]*scope{}, tokens: map[string]string{}, now: time.Now}
}
func randomID() (string, error) {
	b := make([]byte, 24)
	if _, e := rand.Read(b); e != nil {
		return "", e
	}
	return hex.EncodeToString(b), nil
}
func (b *Broker) NewScope(projectID string) (string, string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed || len(b.scopes) >= MaxScopes {
		return "", "", app.Fail("interaction_unavailable", "conversation capacity reached or broker closed")
	}
	if strings.TrimSpace(projectID) == "" {
		return "", "", app.Fail("invalid_arguments", "project scope is required")
	}
	id, e := randomID()
	if e != nil {
		return "", "", e
	}
	token, e := randomID()
	if e != nil {
		return "", "", e
	}
	b.scopes[id] = &scope{project: projectID, token: token}
	b.tokens[token] = id
	return id, token, nil
}
func (b *Broker) ScopeForToken(token string) (string, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	id, ok := b.tokens[token]
	return id, ok && !b.closed
}
func (b *Broker) get(id string) (*scope, error) {
	s := b.scopes[id]
	if s == nil {
		return nil, app.Fail("not_found", "conversation scope is unavailable")
	}
	return s, nil
}
func validText(s string, max int) bool { return utf8.ValidString(s) && len(s) <= max }
func (b *Broker) Ask(id string, in AskInput) (Question, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	s, e := b.get(id)
	if e != nil {
		return Question{}, e
	}
	if s.closed || b.closed {
		return Question{}, app.Fail("interaction_unavailable", "conversation is closed")
	}
	if len(s.questions) >= MaxQuestions {
		return Question{}, app.Fail("interaction_unavailable", "question capacity reached")
	}
	if strings.TrimSpace(in.Prompt) == "" || !validText(in.Prompt, MaxText) {
		return Question{}, app.Fail("invalid_arguments", "prompt is required and bounded to 16 KiB")
	}
	switch in.Kind {
	case "text":
		if len(in.Choices) != 0 {
			return Question{}, app.Fail("invalid_arguments", "text question cannot include choices")
		}
	case "single_choice", "multiple_choice":
		if len(in.Choices) == 0 || len(in.Choices) > 32 {
			return Question{}, app.Fail("invalid_arguments", "choice question needs 1 to 32 choices")
		}
	default:
		return Question{}, app.Fail("invalid_arguments", "unknown response kind")
	}
	seen := map[string]bool{}
	for _, c := range in.Choices {
		if strings.TrimSpace(c.ID) == "" || !validText(c.ID, 128) || strings.TrimSpace(c.Label) == "" || !validText(c.Label, 1024) || seen[c.ID] {
			return Question{}, app.Fail("invalid_arguments", "choice identifiers/labels must be nonempty, bounded and unique")
		}
		seen[c.ID] = true
	}
	if in.PreviousRequestID != "" {
		if _, e = b.find(s, in.PreviousRequestID); e != nil {
			return Question{}, e
		}
	}
	qid, e := randomID()
	if e != nil {
		return Question{}, e
	}
	now := b.now()
	q := Question{ID: qid, ScopeID: id, Prompt: in.Prompt, Kind: in.Kind, Choices: append([]Choice(nil), in.Choices...), PreviousRequestID: in.PreviousRequestID, Status: "pending", CreatedAt: now, ExpiresAt: now.Add(QuestionLifetime)}
	s.questions = append(s.questions, &entry{q: q, done: make(chan struct{})})
	return cloneQ(q), nil
}
func (b *Broker) find(s *scope, id string) (*entry, error) {
	for _, e := range s.questions {
		if e.q.ID == id {
			b.expire(e)
			return e, nil
		}
	}
	return nil, app.Fail("not_found", "request is not owned by this conversation")
}
func (b *Broker) expire(e *entry) {
	if e.q.Status == "pending" && !b.now().Before(e.q.ExpiresAt) {
		e.q.Status = "expired"
		close(e.done)
	}
}
func cloneQ(q Question) Question {
	q.Choices = append([]Choice(nil), q.Choices...)
	q.OptionIDs = append([]string(nil), q.OptionIDs...)
	return q
}
func (b *Broker) Wait(ctx context.Context, id, request string, timeout time.Duration) (Question, error) {
	if timeout < 0 || timeout > 30*time.Second {
		return Question{}, app.Fail("invalid_arguments", "wait must be between zero and 30 seconds")
	}
	b.mu.Lock()
	s, e := b.get(id)
	if e != nil {
		b.mu.Unlock()
		return Question{}, e
	}
	v, e := b.find(s, request)
	if e != nil {
		b.mu.Unlock()
		return Question{}, e
	}
	q := cloneQ(v.q)
	done := v.done
	b.mu.Unlock()
	if q.Status != "pending" || timeout == 0 {
		return q, nil
	}
	remaining := time.Until(q.ExpiresAt)
	if remaining < timeout {
		timeout = remaining
	}
	if timeout < 0 {
		timeout = 0
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return Question{}, ctx.Err()
	case <-done:
	case <-timer.C:
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.expire(v)
	return cloneQ(v.q), nil
}

// Answer is a host-only API. It is deliberately absent from ToolCall and MCP tools.
func (b *Broker) Answer(id, request string, in AnswerInput) (Question, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	s, e := b.get(id)
	if e != nil {
		return Question{}, e
	}
	v, e := b.find(s, request)
	if e != nil {
		return Question{}, e
	}
	if in.Status != "answered" && in.Status != "declined" && in.Status != "cancelled" {
		return Question{}, app.Fail("invalid_arguments", "invalid user disposition")
	}
	if !validText(in.Text, MaxText) {
		return Question{}, app.Fail("invalid_arguments", "answer text is too large or invalid UTF-8")
	}
	if in.Status != "answered" && (in.Text != "" || len(in.OptionIDs) != 0) {
		return Question{}, app.Fail("invalid_arguments", "non-answer disposition cannot contain an answer")
	}
	if in.Status == "answered" {
		if strings.TrimSpace(in.Text) == "" && len(in.OptionIDs) == 0 {
			return Question{}, app.Fail("invalid_arguments", "submit text or a selected choice")
		}
		if v.q.Kind == "text" && len(in.OptionIDs) > 0 || v.q.Kind == "single_choice" && len(in.OptionIDs) > 1 {
			return Question{}, app.Fail("invalid_arguments", "selection does not match response kind")
		}
		valid := map[string]bool{}
		for _, c := range v.q.Choices {
			valid[c.ID] = true
		}
		seen := map[string]bool{}
		for _, x := range in.OptionIDs {
			if !valid[x] || seen[x] {
				return Question{}, app.Fail("invalid_arguments", "unknown or repeated choice")
			}
			seen[x] = true
		}
	}
	if v.q.Status != "pending" {
		if v.q.Status == in.Status && v.q.Text == in.Text && reflect.DeepEqual(v.q.OptionIDs, normalizeIDs(in.OptionIDs)) {
			return cloneQ(v.q), nil
		}
		return Question{}, app.Fail("stale_revision", "request already has a terminal disposition")
	}
	v.q.Status = in.Status
	v.q.Text = in.Text
	v.q.OptionIDs = normalizeIDs(in.OptionIDs)
	close(v.done)
	return cloneQ(v.q), nil
}
func normalizeIDs(v []string) []string {
	if len(v) == 0 {
		return nil
	}
	return append([]string(nil), v...)
}
func (b *Broker) Present(id string, in ViewInput) (View, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	s, e := b.get(id)
	if e != nil {
		return View{}, e
	}
	if b.closed || s.closed || len(s.views) >= MaxViews {
		return View{}, app.Fail("interaction_unavailable", "view capacity reached or conversation closed")
	}
	if strings.TrimSpace(in.Title) == "" || !validText(in.Title, 256) || strings.TrimSpace(in.Fallback) == "" || !validText(in.Fallback, MaxText) || strings.TrimSpace(in.HTML) == "" || !validText(in.HTML, MaxHTML) || len(in.Data) > MaxHTML || len(in.Data) > 0 && !json.Valid(in.Data) {
		return View{}, app.Fail("invalid_arguments", "view requires bounded title, fallback, HTML and valid JSON data")
	}
	if in.RequestID != "" {
		if _, e = b.find(s, in.RequestID); e != nil {
			return View{}, e
		}
	}
	vid, e := randomID()
	if e != nil {
		return View{}, e
	}
	v := View{ID: vid, ResourceURI: "ui://factoryctl/" + vid, Title: in.Title, Fallback: in.Fallback, HTML: in.HTML, Data: append(json.RawMessage(nil), in.Data...), RequestID: in.RequestID}
	s.views = append(s.views, v)
	v.Data = append(json.RawMessage(nil), v.Data...)
	return v, nil
}
func (b *Broker) Snapshot(id string) (Snapshot, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	s, e := b.get(id)
	if e != nil {
		return Snapshot{}, e
	}
	v := Snapshot{ScopeID: id, ProjectID: s.project, Closed: s.closed, Questions: []Question{}, Views: []View{}}
	for _, q := range s.questions {
		b.expire(q)
		v.Questions = append(v.Questions, cloneQ(q.q))
	}
	for _, x := range s.views {
		x.Data = append(json.RawMessage(nil), x.Data...)
		v.Views = append(v.Views, x)
	}
	return v, nil
}

// CancelPending ends outstanding questions for an interrupted turn, without
// revoking the conversation credential or preventing a later follow-up turn.
func (b *Broker) CancelPending(id, disposition string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	s := b.scopes[id]
	if s == nil {
		return
	}
	if disposition != "cancelled" && disposition != "unavailable" {
		disposition = "cancelled"
	}
	for _, q := range s.questions {
		if q.q.Status == "pending" {
			q.q.Status = disposition
			close(q.done)
		}
	}
}

func (b *Broker) CancelScope(id, disposition string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	s := b.scopes[id]
	if s == nil {
		return
	}
	if disposition != "cancelled" && disposition != "unavailable" {
		disposition = "cancelled"
	}
	s.closed = true
	delete(b.tokens, s.token)
	for _, q := range s.questions {
		if q.q.Status == "pending" {
			q.q.Status = disposition
			close(q.done)
		}
	}
}
func (b *Broker) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.closed = true
	b.tokens = map[string]string{}
	for _, s := range b.scopes {
		s.closed = true
		for _, q := range s.questions {
			if q.q.Status == "pending" {
				q.q.Status = "unavailable"
				close(q.done)
			}
		}
	}
}
func decode(data json.RawMessage, target any) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if e := d.Decode(target); e != nil {
		return app.Fail("invalid_arguments", e.Error())
	}
	if e := d.Decode(new(any)); e != io.EOF {
		return app.Fail("invalid_arguments", "expected one JSON object")
	}
	return nil
}

// ToolCall only accepts model-facing operations; host answer and scope selection are not tools.
func (b *Broker) ToolCall(ctx context.Context, scopeID, name string, args json.RawMessage) (any, error) {
	switch name {
	case "factory_user_ask":
		var in AskInput
		if e := decode(args, &in); e != nil {
			return nil, e
		}
		return b.Ask(scopeID, in)
	case "factory_user_wait":
		var in struct {
			RequestID      string `json:"request_id"`
			TimeoutSeconds *int   `json:"timeout_seconds,omitempty"`
		}
		if e := decode(args, &in); e != nil {
			return nil, e
		}
		seconds := 30
		if in.TimeoutSeconds != nil {
			seconds = *in.TimeoutSeconds
		}
		if seconds < 0 || seconds > 30 {
			return nil, app.Fail("invalid_arguments", "wait must be between zero and 30 seconds")
		}
		return b.Wait(ctx, scopeID, in.RequestID, time.Duration(seconds)*time.Second)
	case "factory_ui_present":
		var in ViewInput
		if e := decode(args, &in); e != nil {
			return nil, e
		}
		return b.Present(scopeID, in)
	default:
		return nil, app.Fail("unsupported_capability", fmt.Sprintf("unknown interaction tool %q", name))
	}
}
