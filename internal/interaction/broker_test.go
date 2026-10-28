package interaction

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"
	"time"
)

func scopeFor(t *testing.T, b *Broker) (string, string) {
	t.Helper()
	id, tok, e := b.NewScope("project-instance")
	if e != nil {
		t.Fatal(e)
	}
	return id, tok
}
func askFor(t *testing.T, b *Broker, id string) Question {
	t.Helper()
	q, e := b.Ask(id, AskInput{Prompt: "What should happen?", Kind: "text"})
	if e != nil {
		t.Fatal(e)
	}
	return q
}
func TestPendingNeverBecomesAnAnswer(t *testing.T) {
	b := New()
	id, _ := scopeFor(t, b)
	q := askFor(t, b, id)
	got, e := b.Wait(context.Background(), id, q.ID, time.Millisecond)
	if e != nil || got.Status != "pending" || got.Text != "" {
		t.Fatalf("timeout fabricated answer: %#v %v", got, e)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, e = b.Wait(ctx, id, q.ID, time.Second); e == nil {
		t.Fatal("cancelled wait succeeded")
	}
	got, _ = b.Wait(context.Background(), id, q.ID, 0)
	if got.Status != "pending" {
		t.Fatal("wait cancellation changed question")
	}
	if _, e = b.ToolCall(context.Background(), id, "factory_user_answer", json.RawMessage(`{"text":"invented"}`)); e == nil {
		t.Fatal("model answered itself")
	}
}
func TestQuestionScopesAndHostSubmission(t *testing.T) {
	b := New()
	a, token := scopeFor(t, b)
	other, _ := scopeFor(t, b)
	q := askFor(t, b, a)
	if _, e := b.Wait(context.Background(), other, q.ID, 0); e == nil {
		t.Fatal("cross-scope read")
	}
	if _, e := b.Answer(other, q.ID, AnswerInput{Status: "answered", Text: "spoof"}); e == nil {
		t.Fatal("cross-scope answer")
	}
	if _, e := b.Ask(other, AskInput{Prompt: "Followup", Kind: "text", PreviousRequestID: q.ID}); e == nil {
		t.Fatal("cross-scope followup")
	}
	if id, ok := b.ScopeForToken(token); !ok || id != a {
		t.Fatal("token mapping")
	}
	got, e := b.Answer(a, q.ID, AnswerInput{Status: "answered", Text: " exact user text \n"})
	if e != nil || got.Text != " exact user text \n" {
		t.Fatalf("verbatim answer: %#v %v", got, e)
	}
	if _, e = b.Answer(a, q.ID, AnswerInput{Status: "answered", Text: got.Text}); e != nil {
		t.Fatal("identical repeat must be idempotent", e)
	}
	if _, e = b.Answer(a, q.ID, AnswerInput{Status: "answered", Text: "different"}); e == nil {
		t.Fatal("conflicting repeat")
	}
	follow, e := b.Ask(a, AskInput{Prompt: "Why?", Kind: "text", PreviousRequestID: q.ID})
	if e != nil || follow.PreviousRequestID != q.ID {
		t.Fatal("same-scope followup", e)
	}
}
func TestChoiceValidationAndCustomAnswers(t *testing.T) {
	b := New()
	id, _ := scopeFor(t, b)
	for _, in := range []AskInput{{Prompt: "", Kind: "text"}, {Prompt: "Q", Kind: "other"}, {Prompt: "Q", Kind: "single_choice"}, {Prompt: "Q", Kind: "text", Choices: []Choice{{ID: "a", Label: "A"}}}, {Prompt: "Q", Kind: "single_choice", Choices: []Choice{{ID: "a", Label: "A"}, {ID: "a", Label: "B"}}}} {
		if _, e := b.Ask(id, in); e == nil {
			t.Fatalf("accepted %#v", in)
		}
	}
	q, e := b.Ask(id, AskInput{Prompt: "Choose", Kind: "single_choice", Choices: []Choice{{ID: "a", Label: "A"}, {ID: "b", Label: "B"}}})
	if e != nil {
		t.Fatal(e)
	}
	for _, in := range []AnswerInput{{Status: "answered"}, {Status: "answered", OptionIDs: []string{"missing"}}, {Status: "answered", OptionIDs: []string{"a", "b"}}, {Status: "declined", Text: "answer"}} {
		if _, e = b.Answer(id, q.ID, in); e == nil {
			t.Fatalf("accepted %#v", in)
		}
	}
	got, e := b.Answer(id, q.ID, AnswerInput{Status: "answered", Text: "custom choice"})
	if e != nil || got.Text != "custom choice" {
		t.Fatal("custom choice", e)
	}
	multi, e := b.Ask(id, AskInput{Prompt: "Both?", Kind: "multiple_choice", Choices: []Choice{{ID: "a", Label: "A"}, {ID: "b", Label: "B"}}})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = b.Answer(id, multi.ID, AnswerInput{Status: "answered", OptionIDs: []string{"a", "b"}}); e != nil {
		t.Fatal(e)
	}
}
func TestWaitWakesOnAnswerAndShutdown(t *testing.T) {
	b := New()
	id, tok := scopeFor(t, b)
	q := askFor(t, b, id)
	result := make(chan Question, 1)
	go func() { got, _ := b.Wait(context.Background(), id, q.ID, time.Second); result <- got }()
	b.Answer(id, q.ID, AnswerInput{Status: "declined"})
	if got := <-result; got.Status != "declined" {
		t.Fatal(got)
	}
	q = askFor(t, b, id)
	go func() { got, _ := b.Wait(context.Background(), id, q.ID, time.Second); result <- got }()
	b.Close()
	if got := <-result; got.Status != "unavailable" {
		t.Fatal(got)
	}
	if _, ok := b.ScopeForToken(tok); ok {
		t.Fatal("shutdown token live")
	}
	if _, e := New().Wait(context.Background(), id, q.ID, 0); e == nil {
		t.Fatal("restart reused request")
	}
}
func TestExpiryCancellationAndIsolationCopies(t *testing.T) {
	b := New()
	id, tok := scopeFor(t, b)
	q := askFor(t, b, id)
	b.now = func() time.Time { return q.ExpiresAt.Add(time.Second) }
	got, _ := b.Wait(context.Background(), id, q.ID, 0)
	if got.Status != "expired" {
		t.Fatal(got)
	}
	if _, e := b.Answer(id, q.ID, AnswerInput{Status: "answered", Text: "late"}); e == nil {
		t.Fatal("late answer")
	}
	b.now = time.Now
	q = askFor(t, b, id)
	b.CancelScope(id, "cancelled")
	got, _ = b.Wait(context.Background(), id, q.ID, 0)
	if got.Status != "cancelled" {
		t.Fatal(got)
	}
	if _, ok := b.ScopeForToken(tok); ok {
		t.Fatal("cancelled token live")
	}
	if _, e := b.Ask(id, AskInput{Prompt: "Q", Kind: "text"}); e == nil {
		t.Fatal("closed scope accepted ask")
	}
}
func TestGeneratedViewBoundsAndOwnership(t *testing.T) {
	b := New()
	id, _ := scopeFor(t, b)
	other, _ := scopeFor(t, b)
	q := askFor(t, b, id)
	input := ViewInput{Title: "Chart", Fallback: "Numbers", HTML: "<h1>Generated</h1>", Data: json.RawMessage(`{"a":1}`), RequestID: q.ID}
	if _, e := b.Present(other, input); e == nil {
		t.Fatal("cross-scope view attachment")
	}
	v, e := b.Present(id, input)
	if e != nil {
		t.Fatal(e)
	}
	v.Data[2] = 'X'
	input.Data[2] = 'Y'
	snap, _ := b.Snapshot(id)
	if string(snap.Views[0].Data) != `{"a":1}` {
		t.Fatal("aliased data")
	}
	snap.Views[0].Data[2] = 'Z'
	snap, _ = b.Snapshot(id)
	if string(snap.Views[0].Data) != `{"a":1}` {
		t.Fatal("snapshot aliases data")
	}
	input.HTML = strings.Repeat("x", MaxHTML+1)
	if _, e := b.Present(id, input); e == nil {
		t.Fatal("oversize view")
	}
	if _, e = b.ToolCall(context.Background(), id, "factory_user_ask", json.RawMessage(`{"prompt":"Q","kind":"text","scope_id":"other"}`)); e == nil {
		t.Fatal("model controlled scope")
	}
}
func TestConcurrentAnswersOnlyOneWins(t *testing.T) {
	b := New()
	id, _ := scopeFor(t, b)
	q := askFor(t, b, id)
	var wg sync.WaitGroup
	wins := make(chan string, 2)
	for _, text := range []string{"A", "B"} {
		wg.Add(1)
		go func(text string) {
			defer wg.Done()
			if _, e := b.Answer(id, q.ID, AnswerInput{Status: "answered", Text: text}); e == nil {
				wins <- text
			}
		}(text)
	}
	wg.Wait()
	close(wins)
	if len(wins) != 1 {
		t.Fatal("concurrent answers both applied")
	}
}

func TestInterruptedTurnKeepsConversationAvailable(t *testing.T) {
	b := New()
	id, token := scopeFor(t, b)
	q := askFor(t, b, id)
	b.CancelPending(id, "cancelled")
	got, _ := b.Wait(context.Background(), id, q.ID, 0)
	if got.Status != "cancelled" {
		t.Fatal(got)
	}
	if mapped, ok := b.ScopeForToken(token); !ok || mapped != id {
		t.Fatal("interrupt revoked conversation credential")
	}
	next := askFor(t, b, id)
	if next.Status != "pending" || next.ID == q.ID {
		t.Fatal(next)
	}
}
