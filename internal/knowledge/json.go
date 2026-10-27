package knowledge

import (
	"encoding/json"
	"math"
)

// MarshalJSON keeps generic YAML nonfinite values visible. JSON has no numeric
// representation for them, so only the transport projection uses YAML strings;
// Raw and Frontmatter retain their exact source/semantic values for validation.
func (d Document) MarshalJSON() ([]byte, error) {
	type plain Document
	projected, _ := jsonValue(d.Frontmatter)
	return json.Marshal(struct {
		*plain
		Frontmatter any `json:"frontmatter"`
	}{(*plain)(&d), projected})
}
func jsonValue(v any) (any, bool) {
	switch x := v.(type) {
	case float64:
		if math.IsNaN(x) {
			return ".nan", true
		}
		if math.IsInf(x, 1) {
			return ".inf", true
		}
		if math.IsInf(x, -1) {
			return "-.inf", true
		}
	case map[string]any:
		out := map[string]any{}
		found := false
		for k, item := range x {
			p, b := jsonValue(item)
			out[k] = p
			found = found || b
		}
		return out, found
	case []any:
		out := make([]any, len(x))
		found := false
		for i, item := range x {
			p, b := jsonValue(item)
			out[i] = p
			found = found || b
		}
		return out, found
	}
	return v, false
}
