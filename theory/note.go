package theory

import (
	"fmt"
	"io"
	"math"
	"strings"
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
type Enharmonic int

const (
	EnharmonicSharp Enharmonic = 1
	EnharmonicFlat  Enharmonic = -1
)

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

func TransposeNote(cfg *Config, n Note, degreeClassInterval int, pitchClassInterval int) Note {
	newDegreeClass := adjustDegreeClass(cfg, n.DegreeClass, degreeClassInterval)
	newPitchClass := adjustPitchClass(cfg, n.PitchClass, pitchClassInterval)

	pitchClassDeltaFromDegreeClasses := pitchClassDelta(cfg, pitchClassFromDegreeClass(cfg, n.DegreeClass), pitchClassFromDegreeClass(cfg, newDegreeClass))
	pitchClassDeltaFromPitchClass := pitchClassDelta(cfg, n.PitchClass, newPitchClass)

	accidentalsOffset := pitchClassDeltaFromPitchClass - pitchClassDeltaFromDegreeClasses

	newAccidentals := normalizeAccidentals(cfg, n.Accidentals+accidentalsOffset)

	naturalNoteName := cfg.NaturalNoteNames[newDegreeClass]
	accidentalToken := ""
	if newAccidentals > 0 {
		accidentalToken = strings.Repeat(string(cfg.SharpSymbols[0]), newAccidentals)
	} else if newAccidentals < 0 {
		accidentalToken = strings.Repeat(string(cfg.FlatSymbols[0]), int(math.Abs(float64(newAccidentals))))
	}

	return Note{
		Name:        string(naturalNoteName) + accidentalToken,
		DegreeClass: newDegreeClass,
		PitchClass:  newPitchClass,
		Accidentals: newAccidentals,
	}
}

// func TransposeNote(cfg *Config, n Note, interval int, enharmonic Enharmonic) (Note, error) {
// 	// First, figure out if we should change the degree class. We can do this by looking at the current note's pitch
// 	// class and degree class and see if the bumped pitch class falls into a different degree class.
// 	newPitchClass := n.PitchClass + PitchClass(interval)
// 	newDegreeClass := DegreeClassFromPitchClass(cfg, newPitchClass, enharmonic)

// }

func parseNote(cfg *Config, text string, pos int) (Note, int, error) {
	naturalNoteName, newPos, err := parseNaturalNoteName(cfg, text, pos)
	if err != nil {
		return Note{}, pos, fmt.Errorf("expected natural note name at position %d: %w", newPos, err)
	}

	degreeClass := degreeClassFromNaturalNoteName(cfg, naturalNoteName)
	pitchClass := pitchClassFromDegreeClass(cfg, degreeClass)

	accidentals, newPos := parseAccidentals(cfg, text, newPos)

	return Note{
		Name:        text[pos:newPos],
		DegreeClass: degreeClass,
		PitchClass:  adjustPitchClass(cfg, pitchClass, accidentals),
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
