package parser

import (
	"fmt"
	"strings"
)

// ? enum for parsing outForm normally.
type FormType byte

const (
	Watchable FormType = iota
	// ! Next ones soon.
	JSON
	CSV
)

// ? OutputForm struct.
// ? Provides us to create custom instances to process the data and return it in an expected way
type OutputForm struct {
	Result    *strings.Builder
	Form      FormType
	_taskIter uint16
}

// ? Creates a new instance of the OutputForm struct.
func NewOutputForm(form FormType) *OutputForm {
	return &OutputForm{
		Form:      form,
		Result:    &strings.Builder{},
		_taskIter: 0,
	}
}

// ! WHY???
// TODO: solve.
func (form *OutputForm) Acquire(mask *ParsedMask) error {
	form._taskIter++
	defer form.Release()
	err := form.work(mask)
	return err
}

// TODO pt.2
func (form *OutputForm) Release() {
	form._taskIter--
}

// ! Where to put it?
func (form *OutputForm) Reset() {
	form.Result.Reset()
}

// ! Where to put it?
func (form *OutputForm) String() string {
	return form.Result.String()
}

// ! Local handlers.

func (form *OutputForm) work(mask *ParsedMask) error {
	switch form.Form {
	case Watchable:
		form.Result.WriteString("\n\t\tFound: ")
		form.Result.WriteString(mask.CreateTargetLink())
		form.Result.WriteString("\n")

		// Checking every field
		for k, v := range mask.Fields {
			if k != "" { // Skipping fields with empty names
				form.Result.WriteString("\t- ")
				form.Result.WriteString(k)
				form.Result.WriteString(": ")
				if v.Value != "" {
					form.Result.WriteString(v.Value)
				} else {
					form.Result.WriteString("null")
				}
				form.Result.WriteString("\n")
			}
		}
		return nil
	default:
		return fmt.Errorf("unknown form type")
	}
}
