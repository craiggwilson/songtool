package theory2

import (
	"github.com/craiggwilson/songtool/pkg/theory2/chord"
)

func NameChord(c chord.Chord) string {
	return std.NameChord(c)
}

func ParseChord(text string) (chord.Parsed, error) {
	return std.ParseChord(text)
}
