package theory

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode/utf8"
)

type Chord struct {
	Root      Note
	Intervals []int
	Suffix    string
	Base      Note
}

func (c *Chord) IsValid() bool {
	return c.Root.IsValid()
}

func ParseChord(cfg *Config, text string) (Chord, error) {
	if cfg == nil {
		cfg = &defaultConfig
	}

	root, pos, err := parseNote(cfg, text, 0)
	if err != nil {
		return Chord{}, fmt.Errorf("expected note at position 0: %w", err)
	}

	// Start with major chord intervals.
	intervals := make([]int, len(cfg.MajorChordIntervals))
	copy(intervals, cfg.MajorChordIntervals)

	suffixPos := pos

	// Each modifier group may have multiple applications, but once a group has passed, no more additions.
	for _, modifierGroup := range cfg.ChordModifiers {
		changed := true
		for changed {
			changed = false
			for _, mod := range modifierGroup.Modifiers {
				if strings.HasPrefix(text[pos:], mod.Match) && (len(mod.Except) == 0 || !strings.HasPrefix(text[pos:], mod.Except)) {
					intervals = modifyIntervals(intervals, mod.Intervals)
					pos += len(mod.Match)
					changed = true
				}
			}
		}
	}

	basePos := pos

	base, pos, _ := parseBaseNote(cfg, text, pos)

	if len(text) != pos {
		return Chord{}, fmt.Errorf("expected EOF at position %d, but had %s", pos, text[pos:])
	}

	return Chord{
		Root:      root,
		Intervals: intervals,
		Suffix:    text[suffixPos:basePos],
		Base:      base,
	}, nil
}

func TransposeChord(cfg *Config, chord Chord, degreeClassInterval int, pitchClassInterval int) Chord {
	newRoot := TransposeNoteDirect(cfg, chord.Root, degreeClassInterval, pitchClassInterval)
	newBase := chord.Base
	if chord.Base.IsValid() {
		newBase = TransposeNoteDirect(cfg, chord.Base, degreeClassInterval, pitchClassInterval)
	}

	newIntervals := make([]int, len(chord.Intervals))
	copy(newIntervals, chord.Intervals)

	return Chord{
		Root:      newRoot,
		Intervals: newIntervals,
		Suffix:    chord.Suffix,
		Base:      newBase,
	}
}

func modifyIntervals(intervals []int, modifiers []int) []int {
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

func parseBaseNote(cfg *Config, text string, pos int) (Note, int, error) {
	if len(text) <= pos {
		return Note{}, pos, io.ErrUnexpectedEOF
	}

	v, w := utf8.DecodeRuneInString(text[pos:])
	for _, r := range cfg.BaseNoteDelimiters {
		if v == r {
			return parseNote(cfg, text, pos+w)
		}
	}

	return Note{}, pos, fmt.Errorf("expected one of %q, but got %q", cfg.BaseNoteDelimiters, v)
}
