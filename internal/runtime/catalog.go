package runtime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type Model struct {
	ID            string   `json:"id"`
	Model         string   `json:"model"`
	DisplayName   string   `json:"display_name"`
	Default       bool     `json:"default"`
	DefaultEffort string   `json:"default_effort"`
	Efforts       []string `json:"efforts"`
}

func (c *Client) Models(ctx context.Context) ([]Model, error) { return readModels(ctx, c.Request) }
func readModels(ctx context.Context, request func(context.Context, string, any) (json.RawMessage, error)) ([]Model, error) {
	models := []Model{}
	cursor := ""
	cursors := map[string]bool{}
	aliases := map[string]int{}
	for pageNumber := 0; pageNumber < 20; pageNumber++ {
		params := map[string]any{"limit": 100}
		if cursor != "" {
			params["cursor"] = cursor
		}
		raw, err := request(ctx, "model/list", params)
		if err != nil {
			return nil, fmt.Errorf("cannot read runtime model catalog: %w", err)
		}
		var page struct {
			Data []struct {
				ID            string `json:"id"`
				Model         string `json:"model"`
				DisplayName   string `json:"displayName"`
				Default       bool   `json:"isDefault"`
				DefaultEffort string `json:"defaultReasoningEffort"`
				Efforts       []struct {
					Value string `json:"reasoningEffort"`
				} `json:"supportedReasoningEfforts"`
			} `json:"data"`
			Next *string `json:"nextCursor"`
		}
		if len(raw) > maxMessage || json.Unmarshal(raw, &page) != nil || page.Data == nil || len(page.Data) > 100 {
			return nil, errors.New("invalid runtime model catalog")
		}
		for _, m := range page.Data {
			if m.ID == "" || m.Model == "" || len(m.ID) > 512 || len(m.Model) > 512 || len(m.DisplayName) > 1024 || len(m.Efforts) > 32 {
				return nil, errors.New("invalid runtime model entry")
			}
			item := Model{ID: m.ID, Model: m.Model, DisplayName: m.DisplayName, Default: m.Default, DefaultEffort: m.DefaultEffort, Efforts: []string{}}
			for _, alias := range []string{m.ID, m.Model} {
				if prior, ok := aliases[alias]; ok && prior != len(models) {
					return nil, errors.New("ambiguous runtime model alias")
				}
				aliases[alias] = len(models)
			}
			seen := map[string]bool{}
			for _, e := range m.Efforts {
				if e.Value == "" || len(e.Value) > 64 || seen[e.Value] {
					return nil, errors.New("invalid runtime reasoning choices")
				}
				seen[e.Value] = true
				item.Efforts = append(item.Efforts, e.Value)
			}
			if item.DefaultEffort != "" && !seen[item.DefaultEffort] {
				return nil, errors.New("invalid runtime default reasoning effort")
			}
			models = append(models, item)
		}
		if page.Next == nil || *page.Next == "" {
			if len(models) == 0 {
				return nil, errors.New("empty runtime model catalog")
			}
			return models, nil
		}
		cursor = *page.Next
		if len(cursor) > 4096 || cursors[cursor] {
			return nil, errors.New("invalid runtime model pagination")
		}
		cursors[cursor] = true
	}
	return nil, errors.New("runtime model catalog exceeds page limit")
}
func ValidateChatSelection(s ChatSelection, models []Model) (ChatSelection, error) {
	if err := validChatSelection(s); err != nil {
		return s, err
	}
	if s.Model == "" && s.Effort == "" {
		return s, nil
	}
	var selected *Model
	for i := range models {
		m := &models[i]
		if (s.Model != "" && (s.Model == m.ID || s.Model == m.Model)) || (s.Model == "" && m.Default) {
			if selected != nil {
				return s, errors.New("ambiguous runtime model selection")
			}
			selected = m
		}
	}
	if selected == nil {
		return s, errors.New("selected runtime model or default model could not be verified")
	}
	if s.Effort != "" {
		found := false
		for _, e := range selected.Efforts {
			found = found || e == s.Effort
		}
		if !found {
			return s, errors.New("unsupported reasoning effort for selected runtime model")
		}
	}
	if s.Model != "" {
		s.Model = selected.Model
	}
	return s, nil
}

// ResolveNativeSelection returns a concrete pair, including when resetting to
// native defaults. It never returns or logs unrelated native configuration.
func (c *Client) ResolveNativeSelection(ctx context.Context, desired ChatSelection, models []Model, cwd string) (ChatSelection, error) {
	return resolveNativeSelection(ctx, desired, models, cwd, c.Request)
}
func resolveNativeSelection(ctx context.Context, desired ChatSelection, models []Model, cwd string, request func(context.Context, string, any) (json.RawMessage, error)) (ChatSelection, error) {
	if desired.Model == "" {
		raw, err := request(ctx, "config/read", map[string]any{"cwd": cwd, "includeLayers": false})
		if err != nil {
			return desired, fmt.Errorf("cannot read native model defaults: %w", err)
		}
		var response struct {
			Config *struct {
				Model  *string `json:"model"`
				Effort *string `json:"model_reasoning_effort"`
			} `json:"config"`
		}
		if len(raw) > maxMessage || json.Unmarshal(raw, &response) != nil || response.Config == nil {
			return desired, errors.New("invalid native model defaults")
		}
		if response.Config.Model != nil {
			desired.Model = *response.Config.Model
		}
		if desired.Effort == "" && response.Config.Effort != nil {
			desired.Effort = *response.Config.Effort
		}
		if desired.Model == "" {
			for _, m := range models {
				if m.Default {
					if desired.Model != "" {
						return desired, errors.New("ambiguous native default model")
					}
					desired.Model = m.Model
				}
			}
		}
	}
	if desired.Model == "" {
		return desired, errors.New("native default model unavailable")
	}
	normalized, err := ValidateChatSelection(desired, models)
	if err != nil {
		return desired, err
	}
	if normalized.Effort == "" {
		for _, m := range models {
			if m.Model == normalized.Model {
				normalized.Effort = m.DefaultEffort
				break
			}
		}
	}
	if normalized.Effort == "" {
		return normalized, errors.New("native default reasoning effort unavailable")
	}
	return ValidateChatSelection(normalized, models)
}
