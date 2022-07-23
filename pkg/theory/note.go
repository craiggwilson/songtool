package theory

import (
	"fmt"
	"io"
	"math"
	"strings"
	"unicode/utf8"
)

type Note struct {
	Name        string      `json:"name"`
	DegreeClass DegreeClass `json:"degreeClass"`
	PitchClass  PitchClass  `json:"pitchClass"`
	Accidentals int         `json:"accidentals"`
}

func (n *Note) IsValid() bool {
	return len(n.Name) > 0
}

type DegreeClass int
type PitchClass int
type Enharmonic int

const (
	EnharmonicSharp Enharmonic = 1
	EnharmonicFlat  Enharmonic = -1
)

// func EnhmarmonicFromNote(note Note) Enharmonic {

// }

func MustParseNote(cfg *Config, text string) Note {
	note, err := ParseNote(cfg, text)
	if err != nil {
		panic(err)
	}

	return note
}

func ParseNote(cfg *Config, text string) (Note, error) {
	if cfg == nil {
		cfg = &defaultConfig
	}

	n, pos, err := parseNote(cfg, text, 0)
	if err != nil {
		return Note{}, err
	}

	if len(text) != pos {
		return Note{}, fmt.Errorf("expected EOF at position %d, but had %q", pos, text[pos:])
	}
	return n, err
}

func TransposeNote(cfg *Config, n Note, interval Interval) Note {
	if cfg == nil {
		cfg = &defaultConfig
	}

	newDegreeClass := adjustDegreeClass(cfg, n.DegreeClass, interval.DegreeClass)
	newPitchClass := adjustPitchClass(cfg, n.PitchClass, interval.PitchClass)

	pitchClassDeltaFromDegreeClasses := pitchClassDelta(cfg, pitchClassFromDegreeClass(cfg, n.DegreeClass), pitchClassFromDegreeClass(cfg, newDegreeClass))
	pitchClassDeltaFromPitchClass := pitchClassDelta(cfg, n.PitchClass, newPitchClass)

	newAccidentals := normalizeAccidentals(cfg, n.Accidentals+pitchClassDeltaFromPitchClass-pitchClassDeltaFromDegreeClasses)

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
