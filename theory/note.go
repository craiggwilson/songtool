package theory

import (
	"fmt"
	"io"
	"unicode/utf8"
)

type Note struct {
	Name        string
	DegreeClass DegreeClass
	PitchClass  PitchClass
	Accidentals int
}

type DegreeClass int
type PitchClass int

func CompareNotes(a, b Note) int {
	switch {
	case a.DegreeClass < b.DegreeClass:
		return -1
	case a.DegreeClass > b.DegreeClass:
		return 1
	case a.Accidentals < b.Accidentals:
		return -1
	case a.Accidentals > b.Accidentals:
		return 1
	default:
		return 0
	}
}

func ParseNote(cfg *Config, text string) (Note, error) {
	n, pos, err := parseNote(cfg, text, 0)
	if len(text) != pos {
		return Note{}, fmt.Errorf("expected EOF at position %d, but had %s", pos, text[pos:])
	}
	return n, err
}

func parseNote(cfg *Config, text string, pos int) (Note, int, error) {
	naturalNoteName, newPos, err := parseNaturalNoteName(cfg, text, pos)
	if err != nil {
		return Note{}, pos, fmt.Errorf("expected natural note name at position %d: %w", newPos, err)
	}

	degreeClass := DegreeClassFromNaturalNoteName(cfg, naturalNoteName)
	pitchClass := PitchClassFromDegreeClass(cfg, degreeClass)

	accidentals := 0
	for {
		var accidental int
		accidental, newPos, err = parseAccidental(cfg, text, newPos)
		if err != nil {
			break
		}

		accidentals += accidental
	}

	return Note{
		Name:        text[pos:newPos],
		DegreeClass: degreeClass,
		PitchClass:  AdjustPitchClass(cfg, pitchClass, accidentals),
		Accidentals: accidentals,
	}, newPos, nil
}

func parseNaturalNoteName(cfg *Config, text string, pos int) (rune, int, error) {
	if len(text) <= pos {
		return 0, pos, io.ErrUnexpectedEOF
	}

	v, w := utf8.DecodeRuneInString(text[pos:])

	for _, nn := range cfg.NaturalNoteNames {
		if v == nn {
			return v, pos + w, nil
		}
	}

	return 0, pos, fmt.Errorf("expected natural note name, but got %v", v)
}

func parseAccidental(cfg *Config, text string, pos int) (int, int, error) {
	if len(text) <= pos {
		return 0, pos, io.ErrUnexpectedEOF
	}

	v, w := utf8.DecodeRuneInString(text[pos:])

	for _, ss := range cfg.SharpSymbols {
		if v == ss {
			return 1, pos + w, nil
		}
	}

	for _, fs := range cfg.FlatSymbols {
		if v == fs {
			return -1, pos + w, nil
		}
	}

	return 0, pos, fmt.Errorf("expected sharp or flat, but got %v", v)
}
