package outforms

import (
	"netares/internal/parser"
	json_form "netares/internal/parser/output_forms/json"
	watchable_form "netares/internal/parser/output_forms/watchable"
)

// ? enum for parsing outForm normally.
type FormType byte

const (
	Watchable FormType = iota
	JSON
	// ! soon.
	CSV
)

// ? OutputForm interface.
// ? Allows to use it as an output, just changing the inner logic.

type OutputForm interface {
	FormType() string
	Review(mask *parser.ParsedBody) (string, error)
}

// ? NewOutputForm function.
// ? Creates a new OutputForm object.
// ? It uses factory pattern.
func NewOutputForm(formTypeimpl string) OutputForm {
	switch formTypeimpl {
	case "watchable":
		return watchable_form.NewWatchableForm()
	case "json":
		return json_form.NewJSONForm()
	// case "csv":
	// 	return new(CSVForm)
	default:
		return nil
	}
}
