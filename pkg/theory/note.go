package theory

import (
	"fmt"
	"io"
	"math"
	"strings"
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
type Enharmonic string

const (
	Sharp Enharmonic = "#"
	Flat  Enharmonic = "b"
)

func EnharmonicFromNote(note Note) Enharmonic {
	return defaultTheory.EnharmonicFromNote(note)
}

func (t *Theory) EnharmonicFromNote(note Note) Enharmonic {
	switch {
	case note.Accidentals > 0:
		return Sharp
	case note.Accidentals < 0:
		return Flat
	default:
		// If the note falls after a half-step, it's enharmonic is flat.
		prevDegreeClass := t.Config.AdjustDegreeClass(note.DegreeClass, -1)
		prevDegreeClassPitchClass := t.Config.DegreeClassToPitchClass[prevDegreeClass]
		prevPitchClass := t.Config.AdjustPitchClass(note.PitchClass, -1)
		if prevPitchClass == prevDegreeClassPitchClass {
			return Flat
		}

		return Sharp
	}
}

func MustNote(note Note, err error) Note {
	if err != nil {
		panic(err)
	}

	return note
}

func ParseNote(text string) (Note, error) {
	return defaultTheory.ParseNote(text)
}

func (t *Theory) ParseNote(text string) (Note, error) {
	n, pos, err := t.parseNote(text, 0)
	if err != nil {
		return Note{}, err
	}

	if len(text) != pos {
		return Note{}, fmt.Errorf("expected EOF at position %d, but had %q", pos, text[pos:])
	}
	return n, err
}

func TransposeNote(n Note, interval Interval) Note {
	return defaultTheory.TransposeNote(n, interval)
}

func (t *Theory) TransposeNote(n Note, interval Interval) Note {
	newDegreeClass := t.Config.AdjustDegreeClass(n.DegreeClass, interval.DegreeClass)
	newPitchClass := t.Config.AdjustPitchClass(n.PitchClass, interval.PitchClass)

	pitchClassDeltaFromDegreeClasses := t.Config.PitchClassDelta(t.Config.PitchClassFromDegreeClass(n.DegreeClass), t.Config.PitchClassFromDegreeClass(newDegreeClass))
	pitchClassDeltaFromPitchClass := t.Config.PitchClassDelta(n.PitchClass, newPitchClass)

	newAccidentals := t.Config.NormalizeAccidentals(n.Accidentals + pitchClassDeltaFromPitchClass - pitchClassDeltaFromDegreeClasses)

	naturalNoteName := t.Config.NaturalNoteNames[newDegreeClass]
	accidentalToken := ""
	if newAccidentals > 0 {
		accidentalToken = strings.Repeat(string(t.Config.SharpSymbols[0]), newAccidentals)
	} else if newAccidentals < 0 {
		accidentalToken = strings.Repeat(string(t.Config.FlatSymbols[0]), int(math.Abs(float64(newAccidentals))))
	}

	return Note{
		Name:        string(naturalNoteName) + accidentalToken,
		DegreeClass: newDegreeClass,
		PitchClass:  newPitchClass,
		Accidentals: newAccidentals,
	}
}

func (t *Theory) parseNote(text string, pos int) (Note, int, error) {
	naturalNoteName, newPos, err := t.parseNaturalNoteName(text, pos)
	if err != nil {
		return Note{}, pos, fmt.Errorf("expected natural note name at position %d: %w", newPos, err)
	}

	degreeClass := t.Config.DegreeClassFromNaturalNoteName(naturalNoteName)
	pitchClass := t.Config.PitchClassFromDegreeClass(degreeClass)

	accidentals, newPos := t.parseAccidentals(text, newPos)

	return Note{
		Name:        text[pos:newPos],
		DegreeClass: degreeClass,
		PitchClass:  t.Config.AdjustPitchClass(pitchClass, accidentals),
		Accidentals: accidentals,
	}, newPos, nil
}

func (t *Theory) parseNaturalNoteName(text string, pos int) (string, int, error) {
	if len(text) <= pos {
		return "", pos, io.ErrUnexpectedEOF
	}

	for _, nn := range t.Config.NaturalNoteNames {
		if strings.HasPrefix(text[pos:], nn) {
			return nn, pos + len(nn), nil
		}
	}

	return "", pos, fmt.Errorf("expected one of %q, but got %q", t.Config.NaturalNoteNames, text[pos:])
}

func (t *Theory) parseAccidentals(text string, pos int) (int, int) {
	if len(text) <= pos {
		return 0, pos
	}

	accidentals, pos := t.parseSharps(text, pos)
	if accidentals != 0 {
		return accidentals, pos
	}

	return t.parseFlats(text, pos)
}

func (t *Theory) parseSharps(text string, pos int) (int, int) {
	if len(text) <= pos {
		return 0, pos
	}

	accidentals := 0

	changed := true
	for changed {
		changed = false
		for _, sym := range t.Config.SharpSymbols {
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

func (t *Theory) parseFlats(text string, pos int) (int, int) {
	if len(text) <= pos {
		return 0, pos
	}

	accidentals := 0

	changed := true
	for changed {
		changed = false
		for _, sym := range t.Config.FlatSymbols {
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
