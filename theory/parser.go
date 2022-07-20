package theory

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/craiggwilson/songtools/theory/key"
	"github.com/craiggwilson/songtools/theory/note"
)

func NewParser(cfg Config) *Parser {
	return &Parser{Config: cfg}
}

type Parser struct {
	Config Config

	minorSuffix string
}

// func (p *Parser) ParseChord(text string) (chord.Chord, error) {

// }

func (p *Parser) ParseKey(text string) (key.Key, error) {
	n, pos, err := p.parseNote(text, 0)
	if err != nil {
		return key.Key{}, fmt.Errorf("expected note at position 0: %w", err)
	}

	if len(p.minorSuffix) == 0 {
	Outer:
		for _, comp := range p.Config.ChordComponents {
			for _, entry := range comp.Entries {
				if len(entry.Intervals) == 2 && ((entry.Intervals[0] == 3 && entry.Intervals[1] == -4) || (entry.Intervals[0] == -4 && entry.Intervals[1] == 3)) {
					p.minorSuffix = entry.Match
					break Outer
				}
			}
		}

		if len(p.minorSuffix) == 0 {
			return key.Key{}, fmt.Errorf("unable to determine minor suffix from Config.ChordComponents")
		}
	}

	kind := key.Major
	if strings.HasPrefix(text[pos:], p.minorSuffix) {
		kind = key.Minor
	}

	return key.Key{
		Note: n,
		Kind: kind,
	}, nil
}

func (p *Parser) ParseNote(text string) (note.Note, error) {
	n, _, err := p.parseNote(text, 0)
	return n, err
}

func (p *Parser) parseNote(text string, pos int) (note.Note, int, error) {
	naturalNoteName, newPos, err := p.parseNaturalNoteName(text, 0)
	if err != nil {
		return note.Note{}, pos, fmt.Errorf("expected natural note name at position %d: %w", newPos, err)
	}

	degreeClass, ok := p.Config.DegreeClass(naturalNoteName)
	if !ok {
		return note.Note{}, pos, fmt.Errorf("natural note name %q does not map to a degree class", naturalNoteName)
	}

	pitchClass, ok := p.Config.PitchClassFromDegreeClass(degreeClass)
	if !ok {
		return note.Note{}, pos, fmt.Errorf("degree class %d does not map to a pitch class", degreeClass)
	}

	accidentals := 0
	for {
		var accidental int
		accidental, newPos, err = p.parseSharpOrFlat(text, newPos)
		if err != nil {
			break
		}

		accidentals += accidental
	}

	return note.Note{
		Name:        text[:newPos],
		DegreeClass: degreeClass,
		PitchClass:  pitchClass,
		Accidentals: accidentals,
	}, newPos, nil
}

func (p *Parser) parseNaturalNoteName(text string, pos int) (rune, int, error) {
	if len(text) <= pos {
		return 0, pos, io.ErrUnexpectedEOF
	}

	v, w := utf8.DecodeRuneInString(text[pos:])

	for _, nn := range p.Config.NaturalNoteNames {
		if v == nn {
			return v, pos + w, nil
		}
	}

	return 0, pos, fmt.Errorf("expected natural note name, but got %v", v)
}

func (p *Parser) parseSharpOrFlat(text string, pos int) (int, int, error) {
	if len(text) <= pos {
		return 0, pos, io.ErrUnexpectedEOF
	}

	v, w := utf8.DecodeRuneInString(text[pos:])

	for _, ss := range p.Config.SharpSymbols {
		if v == ss {
			return 1, pos + w, nil
		}
	}

	for _, fs := range p.Config.FlatSymbols {
		if v == fs {
			return -1, pos + w, nil
		}
	}

	return 0, pos, fmt.Errorf("expected sharp or flat, but got %v", v)
}
