package runtime

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

func TestCatalogPaginationAndNormalization(t *testing.T) {
	calls := 0
	models, err := readModels(context.Background(), func(_ context.Context, method string, p any) (json.RawMessage, error) {
		calls++
		if method != "model/list" {
			t.Fatal(method)
		}
		if calls == 1 {
			return json.RawMessage(`{"data":[{"id":"alias","model":"real","displayName":"Real model","isDefault":true,"defaultReasoningEffort":"medium","supportedReasoningEfforts":[{"reasoningEffort":"low"},{"reasoningEffort":"medium"}]}],"nextCursor":"page2"}`), nil
		}
		if p.(map[string]any)["cursor"] != "page2" {
			t.Fatal(p)
		}
		return json.RawMessage(`{"data":[],"nextCursor":null}`), nil
	})
	if err != nil || calls != 2 || len(models) != 1 || models[0].DisplayName != "Real model" {
		t.Fatal(models, err, calls)
	}
	s, err := ValidateChatSelection(ChatSelection{"alias", "low"}, models)
	if err != nil || s.Model != "real" {
		t.Fatal(s, err)
	}
	if _, err = ValidateChatSelection(ChatSelection{"real", "ultra"}, models); err == nil {
		t.Fatal("unsupported effort")
	}
	if _, err = ValidateChatSelection(ChatSelection{"missing", ""}, models); err == nil {
		t.Fatal("unsupported model")
	}
}
func TestCatalogMalformedAndBounds(t *testing.T) {
	for _, raw := range []string{`{}`, `{"data":[]}`, `{"data":[{"id":"a","model":"b"},{"id":"b","model":"c"}]}`, `{"data":[{"id":"a","model":"a","defaultReasoningEffort":"high","supportedReasoningEfforts":[]}]}`, `{"data":[],"nextCursor":"again"}`} {
		t.Run(raw, func(t *testing.T) {
			_, err := readModels(context.Background(), func(context.Context, string, any) (json.RawMessage, error) { return json.RawMessage(raw), nil })
			if err == nil {
				t.Fatal("invalid catalog accepted")
			}
		})
	}
}
func TestNativeDefaultsAndReset(t *testing.T) {
	models := []Model{{ID: "alias", Model: "configured", DefaultEffort: "low", Efforts: []string{"low", "high"}}, {ID: "catalog", Model: "catalog", Default: true, DefaultEffort: "medium", Efforts: []string{"medium"}}}
	tests := []struct {
		name    string
		desired ChatSelection
		config  string
		want    ChatSelection
		calls   int
		fail    bool
	}{
		{"native configured", ChatSelection{}, `{"config":{"model":"alias","model_reasoning_effort":"high"}}`, ChatSelection{"configured", "high"}, 1, false},
		{"explicit resets effort", ChatSelection{"alias", ""}, ``, ChatSelection{"configured", "low"}, 0, false},
		{"effort overrides native", ChatSelection{"", "low"}, `{"config":{"model":"alias","model_reasoning_effort":"high"}}`, ChatSelection{"configured", "low"}, 1, false},
		{"absent native model", ChatSelection{}, `{"config":{"model":null,"model_reasoning_effort":null}}`, ChatSelection{"catalog", "medium"}, 1, false},
		{"unknown native model", ChatSelection{}, `{"config":{"model":"unknown"}}`, ChatSelection{}, 1, true},
		{"malformed", ChatSelection{}, `{"config":null}`, ChatSelection{}, 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			s, err := resolveNativeSelection(context.Background(), tt.desired, models, "/tmp", func(_ context.Context, m string, p any) (json.RawMessage, error) {
				calls++
				if m != "config/read" || p.(map[string]any)["includeLayers"] != false {
					t.Fatal(m, p)
				}
				return json.RawMessage(tt.config), nil
			})
			if calls != tt.calls || (err != nil) != tt.fail || !tt.fail && s != tt.want {
				t.Fatal(s, err, calls)
			}
		})
	}
	_, err := resolveNativeSelection(context.Background(), ChatSelection{}, models, "/tmp", func(context.Context, string, any) (json.RawMessage, error) { return nil, errors.New("offline") })
	if err == nil {
		t.Fatal("config failure hidden")
	}
}

func TestModelCatalogThroughClientRPC(t *testing.T) {
	c, err := helper(t, Config{}, "")
	if err != nil {
		t.Fatal(err)
	}
	models, err := c.Models(context.Background())
	if err != nil || len(models) != 1 || models[0].Model != "qualified" || len(models[0].Efforts) != 1 || models[0].Efforts[0] != "low" {
		t.Fatal(models, err)
	}
}
