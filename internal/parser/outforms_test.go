package parser_test

import (
	"netares/internal/parser"
	"testing"
)

// ? Test for the output form.
func TestPrepareOutput(t *testing.T) {
	pm := new(parser.ParsedMask)
	pm.UnmarshalJSON(testFullMaskExample)
	pm.TargetName = "testName"

	for field := range pm.Fields {
		pm.Fields[field] = parser.ParsedValue{
			Route: pm.Fields[field].Route,
			Value: "field - " + field,
		}
	}

	form := parser.NewOutputForm(parser.Watchable)
	if err := form.Acquire(pm); err != nil {
		t.Fatal(err)
	}

	out := form.String()
	form.Reset()

	t.Log(out)
}
