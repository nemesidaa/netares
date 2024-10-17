package domhttp

import (
	"errors"
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
)

// ? QueryDirective is a struct that contains the query parameters and the query result.
// * It really helps in a many ways, where we need to parse the query down the list.
type QueryDirective struct {
	param  []string
	result string
}

// ? GetParams returns the query parameters.
func (qd *QueryDirective) GetParams() []string {
	return qd.param
}

// ? Result returns the query result.
func (qd *QueryDirective) Result() string {
	return qd.result
}

// ? ParseQuery parses the query and returns the result.
func ParseQuery(route string, body []byte) (*QueryDirective, error) {
	// ? initialization
	qd := &QueryDirective{
		param:  strings.Split(route, ">"),
		result: "",
	}

	if len(qd.param) == 0 {
		return nil, errors.New("invalid query")
	}

	// ? Converting to HTML...
	doc, err := htmlquery.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// ? Handling the query...
	xpathQuery := strings.Join(qd.param, "/")
	node := htmlquery.FindOne(doc, xpathQuery)
	if node == nil {
		return nil, fmt.Errorf("failed to find node for query: %s", xpathQuery)
	}

	qd.result = htmlquery.InnerText(node)
	return qd, nil
}
