package theory2

import (
	"github.com/craiggwilson/songtool/pkg/theory2/chord"
)

type ChordNamer interface {
	NameChord(chord.Chord) string
}

func NameChord(c chord.Chord) string {
	return std.NameChord(c)
}

type ChordParser interface {
	ParseChord(string) (chord.Chord, error)
}

func ParseChord(text string) (chord.Chord, error) {
	return std.ParseChord(text)
}
