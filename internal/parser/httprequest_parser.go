package parser

import (
	"io"
	domhttp "netares/internal/parser/domrouter"
	"strings"
	"time"
)

// ? ParsedValue is actually one of the key second-route values.
// ? It contains all responses, route metadata... I think, that is the containment value,
// * That provides us to be careful with the deeds with her.
type ParsedBody struct {
	target string
	Data   map[string]ParsedValue // ? Format of the map like map[name_val]{route, result}.
	Router *domhttp.DomRouter
}

// ? Creates a new instance of ParsedBody
func NewParsedBody(target string, data map[string]ParsedValue) *ParsedBody {
	routerData := make(map[string]string)
	for key, value := range data {
		routerData[key] = value.Route
	}
	return &ParsedBody{
		target: target,
		Data:   data,
		Router: domhttp.NewRouter(routerData),
	}
}

// ? Changes an object and returns only an error.
func (pb *ParsedBody) Parse(raw io.ReadCloser) error {
	pb.parseBody(raw)
	return nil
}

// ? Parses the body and returns the result.
// ? first logical abstractions are very cute.
func (pb *ParsedBody) parseBody(raw io.ReadCloser) error {

	defer raw.Close()
	result, err := pb.Router.Solve(raw, 5*time.Second)
	if err != nil {
		return err
	}
	for key, value := range result {
		pb.Data[key] = ParsedValue{Route: pb.Data[key].Route, Value: strings.TrimSpace(value)}
	}
	return nil
}

// ? Returns the target name
func (pb *ParsedBody) GetTarget() string {
	return pb.target
}
