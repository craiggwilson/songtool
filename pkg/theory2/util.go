package theory2

import (
	"unicode"
	"unicode/utf8"
)

var degreeClassToPitchClass = [7]int{0, 2, 4, 5, 7, 9, 11}

func consumeWhitespace(text string) int {
	pos := 0
	for _, r := range text {
		if !unicode.IsSpace(r) {
			break
		}

		pos += utf8.RuneLen(r)
	}

	return pos
}
