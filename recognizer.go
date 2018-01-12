package parglare

import (
	"regexp"
	"strings"
)

// Recognizer is a generalization of token recognizer.
type Recognizer interface {
	// Recognize is called with input object, position and custom state
	// which might be used in subsequent calls to recognizers. Should return
	// a token, next position and a new state.
	Recognize(input interface{}, pos int, state interface{}) (
		token interface{}, nextPos int, nextState interface{})
}

// StringRecognizer recognizes tokens by the provided string.
type StringRecognizer struct {
	Name       string
	Value      string
	IgnoreCase bool
}

// RegExRecognizer recognizes tokens by the provided regular expression.
type RegExRecognizer struct {
	Name       string
	RegEx      *regexp.Regexp
	RegExFlags int
	IgnoreCase bool
}

func NewRegExRecognizer(name, regex string, flags int, ignoreCase bool) *RegExRecognizer {
	r := new(RegExRecognizer)
	r.Name = name

	// regex is anchored at the beginning of the string.
	re, error := regexp.Compile("^" + regex)

	if error != nil {
		panic("Invalid regex for recognizer " + name + ". ")
	}

	r.RegEx = re
	r.RegExFlags = flags
	r.IgnoreCase = ignoreCase

	return r
}

func (r *StringRecognizer) Recognize(input interface{}, pos int, state interface{}) (
	token interface{}, nextPos int, nextState interface{}) {

	var in string = input.(string)[pos : pos+len(r.Value)]
	var cmpin string = in
	var cmpval = r.Value
	if r.IgnoreCase {
		cmpin = strings.ToLower(in)
		cmpval = strings.ToLower(cmpval)
	}
	if cmpin == cmpval {
		return in, pos + len(in), state

	} else {
		return "", pos, state
	}
}

func (r *RegExRecognizer) Recognize(input interface{}, pos int, state interface{}) (
	token interface{}, nextPos int, nextState interface{}) {

	var in string = input.(string)[pos:]

	var match = r.RegEx.FindString(in)
	if match == "" {
		return "", pos, nextState
	}

	return match, pos + len(match), state

}
