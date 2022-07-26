package theory2

import (
	"fmt"
	"io"
	"strings"

	"github.com/craiggwilson/songtool/pkg/theory2/interval"
	"github.com/craiggwilson/songtool/pkg/theory2/note"
	"github.com/craiggwilson/songtool/pkg/theory2/scale"
)

func DefaultConfig() *Config {
	return &Config{
		NaturalNoteNames: [7]string{"C", "D", "E", "F", "G", "A", "B"},
		SharpSymbols:     []string{"#"},
		FlatSymbols:      []string{"b"},
		Scales: map[string][]interval.Interval{
			"Major":     interval.Scales.Ionian,
			"Ionian":    interval.Scales.Ionian,
			"Chromatic": interval.Scales.Chromatic,
		},
	}
}

type Config struct {
	NaturalNoteNames [7]string
	SharpSymbols     []string
	FlatSymbols      []string

	Scales map[string][]interval.Interval
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

func (c *Config) ParseNote(text string) (note.Note, error) {
	n, pos, err := c.parseNote(text, 0)
	if err != nil {
		return note.Note{}, err
	}

	if len(text) != pos {
		return note.Note{}, fmt.Errorf("expected EOF at position %d, but had %q", pos, text[pos:])
	}

	return n, err
}

func (c *Config) ParseScale(text string) (scale.Scale, error) {
	parts := strings.SplitN(text, " ", 2)

	root, err := c.ParseNote(parts[0])
	if err != nil {
		return scale.Scale{}, err
	}

	scaleName := "Major"
	if len(parts) == 2 {
		scaleName = parts[1]
	}

	intervals, ok := c.Scales[scaleName]
	if !ok {
		return scale.Scale{}, fmt.Errorf("unknown scale name %q", scaleName)
	}

	return scale.Generate(fmt.Sprintf("%s %s", parts[0], scaleName), root, intervals...), nil
}

func (c *Config) degreeClassFromNaturalNoteName(naturalNoteName string) int {
	for i, nn := range c.NaturalNoteNames {
		if nn == naturalNoteName {
			return i
		}
	}

	panic(fmt.Sprintf("natural note name %q does not map to a degree class", naturalNoteName))
}

func (c *Config) parseAccidentals(text string, pos int) (int, int) {
	if len(text) <= pos {
		return 0, pos
	}

	accidentals, pos := c.parseSharps(text, pos)
	if accidentals != 0 {
		return accidentals, pos
	}

	return c.parseFlats(text, pos)
}

func (c *Config) parseNote(text string, pos int) (note.Note, int, error) {
	naturalNoteName, newPos, err := c.parseNaturalNoteName(text, pos)
	if err != nil {
		return note.Note{}, pos, fmt.Errorf("expected natural note name at position %d: %w", newPos, err)
	}

	degreeClass := c.degreeClassFromNaturalNoteName(naturalNoteName)
	pitchClass := degreeClassToPitchClass[degreeClass]
	accidentals, newPos := c.parseAccidentals(text, newPos)

	return note.New(degreeClass, pitchClass+accidentals), newPos, nil
}

func (c *Config) parseNaturalNoteName(text string, pos int) (string, int, error) {
	if len(text) <= pos {
		return "", pos, io.ErrUnexpectedEOF
	}

	for _, nn := range c.NaturalNoteNames {
		if strings.HasPrefix(text[pos:], nn) {
			return nn, pos + len(nn), nil
		}
	}

	return "", pos, fmt.Errorf("expected one of %q, but got %q", c.NaturalNoteNames, text[pos:])
}

func (c *Config) parseSharps(text string, pos int) (int, int) {
	if len(text) <= pos {
		return 0, pos
	}

	accidentals := 0

	changed := true
	for changed {
		changed = false
		for _, sym := range c.SharpSymbols {
			if strings.HasPrefix(text[pos:], sym) {
				accidentals++
				pos += len(sym)
				changed = true
				break
			}
		}
	}

	return accidentals, pos
}

func (c *Config) parseFlats(text string, pos int) (int, int) {
	if len(text) <= pos {
		return 0, pos
	}

	accidentals := 0

	changed := true
	for changed {
		changed = false
		for _, sym := range c.FlatSymbols {
			if strings.HasPrefix(text[pos:], sym) {
				accidentals--
				pos += len(sym)
				changed = true
				break
			}
		}
	}

	return accidentals, pos
}
