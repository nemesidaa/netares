package selector

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

/*
	! Attention! This is an outdated package to refill the idea of XPath usability.
	? I don`t want to remove it, maybe later I`ll reforge it and make it more usable for consumers.
*/

type matchType uint8

type selMatch struct {
	t matchType
	v string
}

const (
	tag matchType = iota
	id
	class
)

// ? QuerySelector user regex to find a partition of the query string.
// ? Pros of this is the simplicity of the syntax and the ease of use, but in a longstep perspective it is quite complicated to use.
func ParseSelector(param string) (prepMatches []selMatch, err error) {
	re := regexp.MustCompile(`(?P<tag>\w+)(?:\s*(?P<id>#[\w-]+))?(?:\s*(?P<class>\.[\w-]+))?(?:\s*(?P<id>#[\w-]+))?`)
	matches := re.FindStringSubmatch(param)
	if matches == nil {
		return nil, fmt.Errorf("invalid selector: %s", param)
	}
	names := re.SubexpNames()

	// ! NOTE: my task there was to save the order of the matches, and because of this, I have to reverse the order of the matches.
	// ! NOTE: I`m not sure if this is the best way, but it works.
	for i, match := range matches {
		if match != "" {
			switch names[i] {
			case "tag":
				prepMatches = append(prepMatches, selMatch{t: tag, v: match})
			case "id":
				prepMatches = append(prepMatches, selMatch{t: id, v: strings.TrimPrefix(match, "#")})
			case "class":
				prepMatches = append(prepMatches, selMatch{t: class, v: strings.TrimPrefix(match, ".")})
			}
		}
	}

	return prepMatches, nil
}

// ? FindSelection finds the start and end bytearrs of the selection.
// ? Still quiet static because of the regex, but works cool.
func (qs *QuerySelector) FindSelection(param string) (start, end []byte, err error) {
	queuedMatches, err := ParseSelector(param)
	if err != nil {
		return nil, nil, err
	}
	tag := queuedMatches[0].v

	// TODO: strings.Builder, sir!
	startTag := "<" + tag
	endTag := "</" + tag + ">"
	for _, match := range queuedMatches {
		switch match.t {
		case id:
			startTag += fmt.Sprintf(` id="%s"`, match.v)
		case class:
			startTag += fmt.Sprintf(` class="%s"`, match.v)
		}
	}
	startTag += ">"

	start = []byte(startTag)
	end = []byte(endTag)
	return start, end, nil
}

// ? Cut cuts data from the body.
// * This is the last, physical layer that provides us to see, what can we solve using these rules.
// * Logically, it jsut collects indexes of the start and the end, and cuts useless partition of them, leaving only the data we need.
// * Or still not cutted again)))
func (qs *QuerySelector) cutData(param string) ([]byte, error) {
	start, end, err := qs.FindSelection(param)
	if err != nil {
		return nil, err
	}

	startIndex := bytes.Index(qs.Data, start)
	if startIndex == -1 {
		return nil, fmt.Errorf("failed to find start: %s", start)
	}

	startIndex += len(start)

	endIndex := bytes.Index(qs.Data[startIndex:], end)
	if endIndex == -1 {
		return nil, fmt.Errorf("failed to find end: %s", end)
	}
	endIndex += startIndex

	return qs.Data[startIndex:endIndex], nil
}
