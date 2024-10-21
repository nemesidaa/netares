package json_form

import (
	"encoding/json"
	"netares/internal/parser"
	"time"
)

// ? marshall is a function for marshalling body to json.
func (form *JSONForm) marshall(body *parser.ParsedBody) ([]byte, error) {
	raw := make(map[string]any)
	raw["target"] = body.GetTarget()
	raw["time"] = time.Now().Format("15:04:05.000000")
	raw["data"] = make(map[string]any)
	for name, field := range body.Data {
		raw["data"].(map[string]any)[name] = field.Value // ! mapping DATA fields there
	}
	return json.Marshal(raw)
}
