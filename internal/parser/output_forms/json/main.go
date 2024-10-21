package json_form

import (
	"context"
	"netares/internal/parser"
)

// ? JSONForm struct.
// * Implements OutputForm interface
type JSONForm struct {
	Ctx context.Context
}

// ? NewJSONForm creates a new JSONForm instance with data
func NewJSONForm() *JSONForm {
	return &JSONForm{
		Ctx: context.Background(),
	}
}

// ? FormType returns the form type
func (form *JSONForm) FormType() string {
	return "json"
}

// ? Review function, the main of the form
func (form *JSONForm) Review(body *parser.ParsedBody) (string, error) {
	// ! Marshalling here
	resp, err := form.marshall(body)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}
