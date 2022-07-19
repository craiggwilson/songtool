package theory

import (
	"github.com/craiggwilson/songtools/theory/note"
)

const DegreeClassCount = 7
const PitchClassCount = 12

var degreeClassToPitchClass = []note.PitchClass{0, 2, 4, 5, 7, 9, 11}

func DefaultConfig() Config {
	return Config{
		NaturalNoteNames: []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G'},
		SharpSymbols:     []rune{'#'},
		FlatSymbols:      []rune{'b'},
	}
}

type Config struct {
	NaturalNoteNames []rune
	SharpSymbols     []rune
	FlatSymbols      []rune
}

func (c *Config) DegreeClass(naturalNoteName rune) (note.DegreeClass, bool) {
	for i, nn := range c.NaturalNoteNames {
		if nn == naturalNoteName {
			return note.DegreeClass(i), true
		}
	}

	return 0, false
}

func (C *Config) PitchClassFromDegreeClass(degreeClass note.DegreeClass) (note.PitchClass, bool) {
	if int(degreeClass) < len(degreeClassToPitchClass) {
		return degreeClassToPitchClass[int(degreeClass)], true
	}

	return 0, false
}
