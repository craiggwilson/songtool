package theory2

import (
	"strings"

	"github.com/craiggwilson/songtool/pkg/theory2/note"
)

func DefaultConfig() *Config {
	return &Config{
		NaturalNoteNames: [7]string{"C", "D", "E", "F", "G", "A", "B"},
		SharpSymbols:     []string{"#"},
		FlatSymbols:      []string{"b"},
	}
}

type Config struct {
	NaturalNoteNames [7]string
	SharpSymbols     []string
	FlatSymbols      []string
}

func (c *Config) NameNote(n note.Note) string {
	degreeClass := n.DegreeClass()
	pitchClass := degreeClassToPitchClass[degreeClass]
	accidentals := n.PitchClass() - pitchClass

	if accidentals > 6 {
		accidentals -= 12
	} else if accidentals < -6 {
		accidentals += 12
	}

	accidentalStr := ""
	if accidentals > 0 {
		accidentalStr = strings.Repeat(c.SharpSymbols[0], accidentals)
	} else if accidentals < 0 {
		accidentalStr = strings.Repeat(c.FlatSymbols[0], -accidentals)
	}

	natural := c.NaturalNoteNames[degreeClass]
	return natural + accidentalStr
}
