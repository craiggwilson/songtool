package theory

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

type Chord struct {
	Root      Note      `json:"root"`
	Semitones []int     `json:"semitones"`
	Suffix    string    `json:"suffix,omitempty"`
	Base      *BaseNote `json:"base,omitempty"`
}

func (c *Chord) IsMinor() bool {
	for _, st := range c.Semitones {
		if st == 3 {
			return true
		}
	}

	return false
}

func (c *Chord) IsValid() bool {
	return c.Root.IsValid()
}

func (c Chord) MarshalJSON() ([]byte, error) {
	type rawChord Chord
	return json.Marshal(struct {
		Name string `json:"name"`
		rawChord
	}{c.Name(), rawChord(c)})
}

func (c *Chord) Name() string {
	name := c.Root.Name + c.Suffix
	if c.Base != nil {
		name += string(c.Base.Delimiter) + c.Base.Name
	}

	return name
}

type BaseNote struct {
	Note
	Delimiter string `json:"delimiter"`
}

func MustChord(chord Chord, err error) Chord {
	if err != nil {
		panic(err)
	}

	return chord
}

func ParseChord(text string) (Chord, error) {
	return std.ParseChord(text)
}

func (t *Theory) ParseChord(text string) (Chord, error) {
	root, pos, err := t.parseNote(text, 0)
	if err != nil {
		return Chord{}, fmt.Errorf("expected note at position 0: %w", err)
	}

	// Start with major chord intervals.
	intervals := make([]int, len(t.Config.MajorChordIntervals))
	copy(intervals, t.Config.MajorChordIntervals)

	suffixPos := pos

	// Each modifier group may have multiple applications, but once a group has passed, no more additions.
	for _, modifierGroup := range t.Config.ChordModifiers {
		changed := true
		for changed {
			changed = false
			for _, mod := range modifierGroup.Modifiers {
				if strings.HasPrefix(text[pos:], mod.Match) && (len(mod.Except) == 0 || !strings.HasPrefix(text[pos:], mod.Except)) {
					intervals = modifySemitones(intervals, mod.Semitones)
					pos += len(mod.Match)
					changed = true
				}
			}
		}
	}

	basePos := pos

	base, pos, _ := t.parseBaseNote(text, pos)

	if len(text) != pos {
		return Chord{}, fmt.Errorf("expected EOF at position %d, but had %s", pos, text[pos:])
	}

	return Chord{
		Root:      root,
		Semitones: intervals,
		Suffix:    text[suffixPos:basePos],
		Base:      base,
	}, nil
}

func TransposeChord(chord Chord, interval Interval) Chord {
	return std.TransposeChord(chord, interval)
}

func (t *Theory) TransposeChord(chord Chord, interval Interval) Chord {
	newRoot := t.TransposeNote(chord.Root, interval)
	newBase := chord.Base
	if chord.Base != nil {
		newBase = &BaseNote{
			Note:      t.TransposeNote(chord.Base.Note, interval),
			Delimiter: chord.Base.Delimiter,
		}
	}

	newIntervals := make([]int, len(chord.Semitones))
	copy(newIntervals, chord.Semitones)

	return Chord{
		Root:      newRoot,
		Semitones: newIntervals,
		Suffix:    chord.Suffix,
		Base:      newBase,
	}
}

func modifySemitones(intervals []int, modifiers []int) []int {
	for i := 0; i < len(modifiers); i++ {
		if modifiers[i] > 0 {
			found := false
			for j := 0; j < len(intervals); j++ {
				if intervals[j] == modifiers[i] {
					found = true
					break
				}
			}

			if !found {
				intervals = append(intervals, modifiers[i])
			}
		}
		if modifiers[i] < 0 {
			mod := modifiers[i] * -1
			for j := 0; j < len(intervals); j++ {
				if intervals[j] == mod {
					for k := j + 1; k < len(intervals); k++ {
						intervals[k-1] = intervals[k]
					}
					intervals = intervals[:len(intervals)-1]
				}
			}
		}
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i] < intervals[j]
	})

	return intervals
}

func (t *Theory) parseBaseNote(text string, pos int) (*BaseNote, int, error) {
	if len(text) <= pos {
		return nil, pos, io.ErrUnexpectedEOF
	}

	for _, r := range t.Config.BaseNoteDelimiters {
		if strings.HasPrefix(text[pos:], r) {
			note, pos, err := t.parseNote(text, pos+len(r))
			return &BaseNote{
				Note:      note,
				Delimiter: r,
			}, pos, err
		}
	}

	return nil, pos, fmt.Errorf("expected one of %q, but got %q", t.Config.BaseNoteDelimiters, text[pos:])
}
