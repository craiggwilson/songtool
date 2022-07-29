package theory

import (
	"github.com/craiggwilson/songtool/pkg/theory/chord"
)

func NameChord(c chord.Chord) string {
	return std.NameChord(c)
}

func ParseChord(text string) (chord.Named, error) {
	return std.ParseChord(text)
}
