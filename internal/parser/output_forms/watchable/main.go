package watchable_form

import (
	"context"
	"fmt"
	"netares/internal/parser"
	"strings"
)

// ? WatchableForm struct.
// * Implements OutputForm interface
type WatchableForm struct {
	Ctx     context.Context
	builder *strings.Builder
}

// ? NewWatchableForm creates a new WatchableForm instance with data
func NewWatchableForm() *WatchableForm {
	return &WatchableForm{
		Ctx:     context.Background(),
		builder: &strings.Builder{},
	}
}

// ? FormType returns the form type
func (form *WatchableForm) FormType() string {
	return "watchable"
}

// ? Review function, the main of the form
func (form *WatchableForm) Review(body *parser.ParsedBody) (string, error) {
	// ! Head using builder.
	form.builder.WriteString("\tRewiew for ")
	form.builder.WriteString(body.GetTarget())
	form.builder.WriteString(" - ")
	form.builder.WriteString(form.FormType())

	for name, field := range body.Data {
		if field.Value == "" {
			fmt.Printf("Warning: field %s is empty, skipping...\n", name)
			continue
		}
		form.builder.WriteString("\n")
		form.builder.WriteString(name)
		form.builder.WriteString(": ")
		form.builder.WriteString(field.Value)
		form.builder.WriteString(";")
	} // ! end of the loop

	return form.builder.String(), nil
}
