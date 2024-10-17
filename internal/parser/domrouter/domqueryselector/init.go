package selector

import (
	"fmt"
)

// ? QuerySelector cuts data from the body.
// ? He has only one field, and works as the independent, using the recursion.
// ? Reaching the finals, it returns the data as a string, in a row saving a plenty of memory(1).
type QuerySelector struct {
	Data []byte // * buffer, at the end it must contain the reached values
}

// ? NewQuerySelector creates a new QuerySelector.
func NewQuerySelector() *QuerySelector {
	return &QuerySelector{}
}

// ? Cut cuts data from the body.
func (qs *QuerySelector) Cut(param []string, body []byte) (string, error) {
	if qs.Data == nil {
		qs.Data = body
	}
	if len(param) != 0 {
		var err error
		// ? logic abstraction, to decompose code.
		qs.Data, err = qs.cutData(param[0])
		if err != nil {
			return "", fmt.Errorf("failed to cut data: %v", err)
		}
		// ? conditional recursion.
		return qs.Cut(param[1:], nil)
	}

	return string(qs.Data), nil
}
