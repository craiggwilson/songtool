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

func ParseNote(cfg *Config, text string) (Note, error) {
	n, pos, err := parseNote(cfg, text, 0)
	if err != nil {
		return Note{}, err
	}

	if len(text) != pos {
		return Note{}, fmt.Errorf("expected EOF at position %d, but had %q", pos, text[pos:])
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

	accidentals, newPos := parseAccidentals(cfg, text, newPos)

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

	return 0, pos, fmt.Errorf("expected one of %q, but got %q", cfg.NaturalNoteNames, v)
}

func parseAccidentals(cfg *Config, text string, pos int) (int, int) {
	if len(text) <= pos {
		return 0, pos
	}

	accidentals, pos := parseSharps(cfg, text, pos)
	if accidentals != 0 {
		return accidentals, pos
	}

	return parseFlats(cfg, text, pos)
}

func parseSharps(cfg *Config, text string, pos int) (int, int) {
	if len(text) <= pos {
		return 0, pos
	}

	accidentals := 0

	changed := true
	v, w := utf8.DecodeRuneInString(text[pos:])
	for changed {
		changed = false
		for _, ss := range cfg.SharpSymbols {
			if v == ss {
				accidentals++
				pos += w
				changed = true
				v, w = utf8.DecodeRuneInString(text[pos:])
				break
			}
		}
	}

	return accidentals, pos
}

func parseFlats(cfg *Config, text string, pos int) (int, int) {
	if len(text) <= pos {
		return 0, pos
	}

	accidentals := 0

	changed := true
	v, w := utf8.DecodeRuneInString(text[pos:])
	for changed {
		changed = false
		for _, ss := range cfg.FlatSymbols {
			if v == ss {
				accidentals--
				pos += w
				changed = true
				v, w = utf8.DecodeRuneInString(text[pos:])
				break
			}
		}
	}

	return accidentals, pos
}
