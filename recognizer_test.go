package parglare

import "testing"

var stringRecognizerTestData = []struct {
	input           string
	pos             int
	ignoreCase      bool
	searchToken     string
	expectedToken   string
	expectedNextPos int
}{
	{"Go Gopher Go!", 3, true, "gopher", "Gopher", 9},
	{"Go Gopher Go!", 3, false, "gopher", "", 3},
	{"Go Gopher Go!", 0, false, "gopher", "", 0},
}

func TestStringRecognizer(t *testing.T) {

	var r Recognizer

	for _, td := range stringRecognizerTestData {

		r = &StringRecognizer{
			td.expectedToken,
			td.expectedToken,
			td.ignoreCase}

		token, nextPos, _ := r.Recognize(td.input, td.pos, nil)

		if token != td.expectedToken || nextPos != td.expectedNextPos {
			t.Errorf("Expected token=%v, nextPos=%v   "+
				"got token=%v, nextPos=%v\n",
				td.expectedToken, td.expectedNextPos,
				token, nextPos)
		}

	}

}

var regexRecognizerTestData = []struct {
	input           string
	pos             int
	searchRegEx     string
	expectedToken   string
	expectedNextPos int
}{
	{"Go 123!", 3, `\d+`, "123", 6},
	{"Go abc!", 3, `\d+`, "", 3},
}

func TestRegExRecognizer(t *testing.T) {

	var r Recognizer

	for _, td := range regexRecognizerTestData {

		r = NewRegExRecognizer("test", td.searchRegEx, 0, false)

		token, nextPos, _ := r.Recognize(td.input, td.pos, nil)

		if token != td.expectedToken || nextPos != td.expectedNextPos {
			t.Errorf("Expected token=%v, nextPos=%v   "+
				"got token=%v, nextPos=%v\n",
				td.expectedToken, td.expectedNextPos,
				token, nextPos)
		}

	}

}
