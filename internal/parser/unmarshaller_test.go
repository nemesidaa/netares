package parser_test

import (
	"netares/internal/parser"
	"testing"
)

var (
	rawTestJSON = []byte(`{
    "source": "https://example.com/*",
    "target": "TargetSystem",
    "data": {
        "field1": {
            "route" : "sus"
        },
        "field2":  {
            "route" : "amogus"
        }
    }
}`)  // ? json for the tests.

	testFullMaskExample = []byte(`{
    "name": "github",
    "source": "https://github.com/*",
    "data": {
      "username": {
        "route": "body#user-profile>div.container>div.vcard-names-container>h1.vcard-names"
      },
      "location": {
        "route": "body#user-profile>div.container>div.vcard-details>li[itemprop='homeLocation']"
      },
      "bio": {
        "route": "body#user-profile>div.container>div.vcard-details>div.user-profile-bio"
      },
      "repositories_count": {
        "route": "body#user-profile>div.container>nav.UnderlineNav>span.Counter"
      },
      "followers_count": {
        "route": "body#user-profile>div.container>div#followers>span.text-bold"
      },
      "following_count": {
        "route": "body#user-profile>div.container>div#following>span.text-bold"
      },
      "profile_picture": {
        "route": "body#user-profile>div.container>img.avatar-user"
      }
    }
  }
  `)  // ? Mask, we need to produce faster, but think, to .
)

// ? Unmarshalling tests.
func TestUnmarshalJSON(t *testing.T) {
	pm := new(parser.ParsedMask)
	err := pm.UnmarshalJSON(rawTestJSON)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	expectedResult := parser.ParsedMask{
		SourceLink: "https://example.com/*",
		Fields: map[string]parser.ParsedValue{
			"field1": {Route: "sus"},
			"field2": {Route: "amogus"},
		},
	}

	if pm.SourceLink != expectedResult.SourceLink {
		t.Fatalf("got %v, want %v", pm.SourceLink, expectedResult.SourceLink)
	}
	if len(pm.Fields) != len(expectedResult.Fields) {
		t.Fatalf("fields mismatch: got %v, want %v", len(pm.Fields), len(expectedResult.Fields))
	}

	t.Logf("Successfully parsed data: %+v", pm)
}

// ? Target link creation link tests.
func TestCreateTargetLink(t *testing.T) {
	pm := parser.ParsedMask{
		SourceLink: "https://example.com/*",
		TargetName: "TargetSystem",
	}

	expectedLink := "https://example.com/TargetSystem"
	if pm.CreateTargetLink() != expectedLink {
		t.Errorf("got %s, wanted %s", pm.CreateTargetLink(), expectedLink)
	} else {
		t.Logf("Target link created correctly: %s", pm.CreateTargetLink())
	}
}

// ? Full mask parsing tests. Like fullmoduled.
func TestFullMaskParsing(t *testing.T) {
	pm := new(parser.ParsedMask)
	err := pm.UnmarshalJSON(testFullMaskExample)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if pm.SourceLink != "https://github.com/*" {
		t.Errorf("expected source link %s, got %s", "https://github.com/*", pm.SourceLink)
	}
	if _, ok := pm.Fields["username"]; !ok {
		t.Errorf("field 'username' not found in mask")
	}

	t.Logf("Successfully parsed mask: %+v", pm)
}
