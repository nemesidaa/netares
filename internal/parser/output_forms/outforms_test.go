package outforms_test

import (
	"netares/internal/parser"
	outforms "netares/internal/parser/output_forms"
	"testing"
)

var (
	customPBody = parser.NewParsedBody(
		"custom",
		map[string]parser.ParsedValue{
			"full-name": {Route: "//body[@id='profile-page']//ul[@id='personal-info-list']//li[@id='full-name']", Value: "John Doe"},
			"username":  {Route: "//body[@id='profile-page']//ul[@id='personal-info-list']//li[@id='username']", Value: "johndoe123"},
		})
)

func TestJSONForm(t *testing.T) {
	json := outforms.NewOutputForm("json")
	resp, err := json.Review(customPBody)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func TestVerboseForm(t *testing.T) {
	verbose := outforms.NewOutputForm("watchable")
	resp, err := verbose.Review(customPBody)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
