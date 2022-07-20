package theory

import (
	"fmt"
	"sort"
	"unicode/utf8"
)

type Key struct {
	Note Note
	Kind KeyKind
}

type KeyKind int

const (
	KeyMajor KeyKind = iota + 1
	KeyMinor
)

func (kk KeyKind) String() string {
	switch kk {
	case KeyMajor:
		return "major"
	case KeyMinor:
		return "minor"
	default:
		return "undefined"
	}
}

func GenerateKeys(cfg *Config, kind KeyKind) []Key {
	keys := make([]Key, 0, len(cfg.NaturalNoteNames)*3)

	for i, nnn := range cfg.NaturalNoteNames {

		degreeClass := DegreeClass(i)
		pitchClass := pitchClassFromDegreeClass(cfg, degreeClass)

		keys = append(keys, Key{
			Note: Note{
				Name:        string(nnn),
				DegreeClass: degreeClass,
				PitchClass:  pitchClass,
				Accidentals: 0,
			},
			Kind: kind,
		})

		keys = append(keys, Key{
			Note: Note{
				Name:        string(nnn) + string(cfg.SharpSymbols[0]),
				DegreeClass: degreeClass,
				PitchClass:  adjustPitchClass(cfg, pitchClass, 1),
				Accidentals: 1,
			},
			Kind: kind,
		})

		keys = append(keys, Key{
			Note: Note{
				Name:        string(nnn) + string(cfg.FlatSymbols[0]),
				DegreeClass: degreeClass,
				PitchClass:  adjustPitchClass(cfg, pitchClass, -1),
				Accidentals: -1,
			},
			Kind: kind,
		})
	}

	SortKeys(keys)

	return keys
}

func ParseKey(cfg *Config, text string) (Key, error) {
	n, pos, err := parseNote(cfg, text, 0)
	if err != nil {
		return Key{}, fmt.Errorf("expected note at position 0: %w", err)
	}

	kind := KeyMajor

	if len(text) > pos {
		v, w := utf8.DecodeRuneInString(text[pos:])
		for _, r := range cfg.MinorKeySymbols {
			if v == r {
				kind = KeyMinor
				pos += w
				break
			}
		}
	}

	if len(text) != pos {
		return Key{}, fmt.Errorf("expected EOF at position %d, but had %q", pos, text[pos:])
	}

	return Key{
		Note: n,
		Kind: kind,
	}, nil
}

func SortKeys(keys []Key) {
	sort.Slice(keys, func(i, j int) bool {
		switch {
		case keys[i].Note.DegreeClass < keys[j].Note.DegreeClass:
			return true
		case keys[i].Note.DegreeClass > keys[j].Note.DegreeClass:
			return false
		case keys[i].Note.Accidentals < keys[j].Note.Accidentals:
			return true
		case keys[i].Note.Accidentals > keys[j].Note.Accidentals:
			return false
		default:
			return false
		}
	})
}
