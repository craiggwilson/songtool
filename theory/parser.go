package theory

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/craiggwilson/songtools/theory/chord"
	"github.com/craiggwilson/songtools/theory/key"
	"github.com/craiggwilson/songtools/theory/note"
)

func NewParser(cfg Config) *Parser {
	return &Parser{Config: cfg}
}

type Parser struct {
	Config Config
}

func (p *Parser) ParseChord(text string) (chord.Chord, error) {
	root, pos, err := p.parseNote(text, 0)
	if err != nil {
		return chord.Chord{}, fmt.Errorf("expected note at position 0: %w", err)
	}

	// Start with major chord intervals.
	intervals := make([]int, len(p.Config.MajorChordIntervals))
	copy(intervals, p.Config.MajorChordIntervals)

	suffixPos := pos

	// Each modifier group may have multiple applications, but once a group has passed, no more additions.
	for _, modifierGroup := range p.Config.ChordModifiers {
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

	base, pos, _ := p.parseBaseNote(text, pos)

	if len(text) != pos {
		return chord.Chord{}, fmt.Errorf("expected EOF at position %d, but had %s", pos, text[pos:])
	}

	return chord.Chord{
		Root:      root,
		Intervals: intervals,
		Suffix:    text[suffixPos:basePos],
		Base:      base,
	}, nil
}

func (p *Parser) ParseKey(text string) (key.Key, error) {
	n, pos, err := p.parseNote(text, 0)
	if err != nil {
		return key.Key{}, fmt.Errorf("expected note at position 0: %w", err)
	}

	kind := key.Major

	if len(text) > pos {
		v, w := utf8.DecodeRuneInString(text[pos:])
		for _, r := range p.Config.MinorKeySymbols {
			if v == r {
				kind = key.Minor
				pos += w
				break
			}
		}
	}

	if len(text) != pos {
		return key.Key{}, fmt.Errorf("expected EOF at position %d, but had %s", pos, text[pos:])
	}

	return key.Key{
		Note: n,
		Kind: kind,
	}, nil
}

func (p *Parser) ParseNote(text string) (note.Note, error) {
	n, pos, err := p.parseNote(text, 0)
	if len(text) != pos {
		return note.Note{}, fmt.Errorf("expected EOF at position %d, but had %s", pos, text[pos:])
	}
	return n, err
}

func (p *Parser) parseBaseNote(text string, pos int) (note.Note, int, error) {
	if len(text) <= pos {
		return note.Note{}, pos, io.ErrUnexpectedEOF
	}

	v, w := utf8.DecodeRuneInString(text[pos:])
	for _, r := range p.Config.BaseNoteDelimiters {
		if v == r {
			return p.parseNote(text, pos+w)
		}
	}

	return note.Note{}, pos, fmt.Errorf("expected base note separator, but got %v", v)
}

func (p *Parser) parseNote(text string, pos int) (note.Note, int, error) {
	naturalNoteName, newPos, err := p.parseNaturalNoteName(text, pos)
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
		Name:        text[pos:newPos],
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
