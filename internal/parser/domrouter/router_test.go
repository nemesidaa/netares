package domhttp_test

import (
	"bytes"
	domhttp "netares/internal/parser/domrouter"
	"testing"
	"time"
)

var rawHtmlBody = []byte(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Test Page</title>
</head>
<body id="profile-page">
    <div id="profile-container">
        <section id="personal-info-section">
            <ul id="personal-info-list">
                <li id="full-name">John Doe</li>
                <li id="username">johndoe123</li>
            </ul>
        </section>
    </div>
</body>
</html>`) // ? The body of the html we want to check on the route modal tests.

// * NO COMMENTS.
type MockReadCloser struct {
	data   []byte
	reader *bytes.Reader
}

func NewMockReadCloser(data []byte) *MockReadCloser {
	return &MockReadCloser{
		data:   data,
		reader: bytes.NewReader(data),
	}
}

func (m *MockReadCloser) Read(p []byte) (n int, err error) {
	return m.reader.Read(p)
}

func (m *MockReadCloser) Close() error {
	return nil
}

// ? XPath module testing.
func TestXPathRouting(t *testing.T) {
	qs := map[string]string{
		"full-name": "//body[@id='profile-page']//ul[@id='personal-info-list']//li[@id='full-name']",
		"username":  "//body[@id='profile-page']//ul[@id='personal-info-list']//li[@id='username']",
	}
	router := domhttp.NewRouter(qs)
	data, err := router.Solve(NewMockReadCloser(rawHtmlBody), 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	for key, value := range data {
		t.Logf("%s: %s", key, value)
	}
}
