package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/tclasen/factoryctl/internal/project"
	rt "github.com/tclasen/factoryctl/internal/runtime"
)

type catalogSnapshot struct {
	Models        []rt.Model
	NativeDefault rt.ChatSelection
	fetched       time.Time
}
type catalogCache struct {
	sync.Mutex
	entries map[string]catalogSnapshot
}
type selectionInput struct {
	Model  *string `json:"model"`
	Effort *string `json:"reasoning_effort"`
}

func (s *selectionInput) selection() (*rt.ChatSelection, error) {
	if s == nil {
		return nil, nil
	}
	if s.Model == nil || s.Effort == nil {
		return nil, fmt.Errorf("selection requires model and reasoning_effort; empty strings select native defaults")
	}
	return &rt.ChatSelection{Model: *s.Model, Effort: *s.Effort}, nil
}
func (c *conversationHost) catalog(ctx context.Context, p project.Project, refresh bool) (catalogSnapshot, error) {
	if c.catalogLookup != nil {
		return c.catalogLookup(ctx, p)
	}
	c.catalogs.Lock()
	defer c.catalogs.Unlock()
	if existing, ok := c.catalogs.entries[p.Root]; ok && !refresh && time.Since(existing.fetched) < 2*time.Minute {
		return existing, nil
	}
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	stop := context.AfterFunc(c.ctx, cancel)
	defer stop()
	client, e := rt.Start(ctx, rt.Config{Executable: c.config.Executable, Dir: p.Root})
	if e != nil {
		return catalogSnapshot{}, e
	}
	defer client.Close()
	models, e := client.Models(ctx)
	if e != nil {
		return catalogSnapshot{}, e
	}
	native, e := client.ResolveNativeSelection(ctx, rt.ChatSelection{}, models, p.Root)
	if e != nil {
		return catalogSnapshot{}, e
	}
	snapshot := catalogSnapshot{Models: models, NativeDefault: native, fetched: time.Now()}
	if c.catalogs.entries == nil {
		c.catalogs.entries = map[string]catalogSnapshot{}
	}
	c.catalogs.entries[p.Root] = snapshot
	return snapshot, nil
}

// Validate against the actual native default, which need not equal model/list's default.
func validateCatalogSelection(desired rt.ChatSelection, catalog catalogSnapshot) (rt.ChatSelection, error) {
	checking := desired
	if checking.Model == "" {
		checking.Model = catalog.NativeDefault.Model
		if checking.Effort == "" {
			checking.Effort = catalog.NativeDefault.Effort
		}
	}
	normalized, e := rt.ValidateChatSelection(checking, catalog.Models)
	if e != nil {
		return desired, e
	}
	if desired.Model != "" {
		desired.Model = normalized.Model
	}
	return desired, nil
}
func (c *conversationHost) selectedDefaults(root string, explicit *rt.ChatSelection) (rt.ChatSelection, error) {
	if explicit != nil {
		return *explicit, nil
	}
	store, e := rt.NewChatSettingsStore(c.config.SettingsPath)
	if e != nil {
		return rt.ChatSelection{}, e
	}
	saved, e := store.Read()
	if e != nil {
		return rt.ChatSelection{}, e
	}
	canonical, e := rt.CanonicalProjectRoot(root)
	if e != nil {
		return rt.ChatSelection{}, e
	}
	return rt.ResolveChatSelection(nil, rt.ChatSelection{Model: c.config.Model, Effort: c.config.Effort}, saved, canonical, c.config.Legacy), nil
}
func (c *conversationHost) settingsHTTP(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("project")
	var input struct {
		Project   string          `json:"project"`
		Scope     string          `json:"scope"`
		Selection json.RawMessage `json:"selection"`
		Revision  string          `json:"revision"`
	}
	if r.Method == "POST" {
		if e := decodeRequest(w, r, &input); e != nil {
			apiError(w, 400, e)
			return
		}
		if projectID != "" && projectID != input.Project {
			apiError(w, 400, fmt.Errorf("Project selection mismatch"))
			return
		}
		projectID = input.Project
		if input.Scope != "global" && input.Scope != "project" {
			apiError(w, 400, fmt.Errorf("Choose global or project defaults"))
			return
		}
		if input.Revision == "" {
			apiError(w, 400, fmt.Errorf("Reload settings before saving"))
			return
		}
	}
	var selected *project.Project
	for i := range c.projects {
		if c.projects[i].Checkout == projectID {
			selected = &c.projects[i]
			break
		}
	}
	if selected == nil {
		apiError(w, 400, fmt.Errorf("Select a loaded project"))
		return
	}
	canonical, e := rt.CanonicalProjectRoot(selected.Root)
	if e != nil {
		apiError(w, 400, e)
		return
	}
	catalog, e := c.catalog(r.Context(), *selected, r.URL.Query().Get("refresh") == "1")
	if e != nil {
		apiError(w, 502, e)
		return
	}
	store, e := rt.NewChatSettingsStore(c.config.SettingsPath)
	if e != nil {
		apiError(w, 500, e)
		return
	}
	saved, readError := store.Read()
	if r.Method == "POST" {
		if readError != nil {
			apiError(w, 400, readError)
			return
		}
		if len(input.Selection) == 0 {
			apiError(w, 400, fmt.Errorf("Explicit selection or null is required to save/reset defaults"))
			return
		}
		var choice *selectionInput
		decoder := json.NewDecoder(bytes.NewReader(input.Selection))
		decoder.DisallowUnknownFields()
		if e := decoder.Decode(&choice); e != nil {
			apiError(w, 400, e)
			return
		}
		selection, e := choice.selection()
		if e != nil {
			apiError(w, 400, e)
			return
		}
		if selection != nil {
			normalized, e := validateCatalogSelection(*selection, catalog)
			if e != nil {
				apiError(w, 400, e)
				return
			}
			selection = &normalized
		}
		scope := ""
		if input.Scope == "project" {
			scope = canonical
		}
		saved, e = store.Save(r.Context(), scope, selection, input.Revision)
		if e != nil {
			code := 400
			if errors.Is(e, rt.ErrSettingsConflict) {
				code = 409
			}
			apiError(w, code, e)
			return
		}
	}
	var override *rt.ChatSelection
	source := "native"
	if c.config.Legacy.Model != "" || c.config.Legacy.Effort != "" {
		source = "personal preferences"
	}
	if saved.Global != nil {
		source = "global"
	}
	if value, ok := saved.Projects[canonical]; ok {
		override = &value
		source = "project"
	}
	if c.config.Model != "" || c.config.Effort != "" {
		source = "invocation"
	}
	effective := rt.ResolveChatSelection(nil, rt.ChatSelection{Model: c.config.Model, Effort: c.config.Effort}, saved, canonical, c.config.Legacy)
	models := []map[string]any{}
	for _, m := range catalog.Models {
		models = append(models, map[string]any{"id": m.ID, "model": m.Model, "display_name": m.DisplayName, "is_default": m.Default, "default_reasoning_effort": m.DefaultEffort, "reasoning_efforts": m.Efforts})
	}
	result := map[string]any{"models": models, "native_default": catalog.NativeDefault, "global": saved.Global, "project": override, "effective": effective, "source": source, "revision": saved.Revision, "config_path": store.Path()}
	if readError != nil {
		result["settings_error"] = readError.Error()
		result["source"] = "settings unavailable"
	}
	apiJSON(w, 200, result)
}
