package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	rawMap                  = map[string]any{}
	errInsuficientReqFields = errors.New("insufficient required fields")
)

// * ParsedMask struct.
// * Parametrized creation, to use an instance, should use.
// * str := new(ParsedMask).
// * Nexly str.Decode(rawDataMask).
type ParsedMask struct {
	SourceLink string                 `json:"source"`
	TargetName string                 `json:"-"`
	Fields     map[string]ParsedValue `json:"data"` // ! As the key it will be name of the searching field, as the value is already described.
}

type ParsedValue struct {
	Route string `json:"route"`
	Value string `json:"-"` // ? will be soon.
}

// * Unmarshalling rules for JSON.
// * e.g. "source": https://smth/* < where * is a wildcard to print TargetName.
// * Data which we need to find is not required, and the format of the outcoming data in the code looks like map[string]any.
func (m *ParsedMask) UnmarshalJSON(raw []byte) error {
	rawMap = make(map[string]any)
	if err := json.Unmarshal(raw, &rawMap); err != nil {
		return fmt.Errorf("failed to unmarshal: %v", err)
	}

	if source, ok := rawMap["source"].(string); ok {
		m.SourceLink = source
	} else {
		return fmt.Errorf("failed to unmarshal: %v", errInsuficientReqFields)
	}

	m.Fields = make(map[string]ParsedValue)

	if dataMap, ok := rawMap["data"].(map[string]interface{}); ok {
		for key, value := range dataMap {
			if fieldMap, ok := value.(map[string]interface{}); ok {
				parsedValue := ParsedValue{}

				if route, ok := fieldMap["route"].(string); ok {
					parsedValue.Route = route
					parsedValue.Value = ""
				} else {
					return fmt.Errorf("missing or invalid 'route' field in data.%s", key)
				}

				m.Fields[key] = parsedValue
			} else {
				return fmt.Errorf("invalid format for data.%s", key)
			}
		}
	} else {
		return fmt.Errorf("missing or invalid 'data' field")
	}

	return nil
}

// ? Creates a link for the target system. Format source.*>source.target.
func (m *ParsedMask) CreateTargetLink() string {
	return strings.ReplaceAll(m.SourceLink, "*", m.TargetName)
}
