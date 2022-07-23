package theory

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Key struct {
	Note   Note    `json:"note"`
	Suffix string  `json:"suffix,omitempty"`
	Kind   KeyKind `json:"kind"`
}

func (k *Key) IsValid() bool {
	return k.Note.IsValid()
}

func (k *Key) MarshalJSON() ([]byte, error) {
	type rawKey Key
	return json.Marshal(struct {
		Name string `json:"name"`
		rawKey
	}{k.Name(), rawKey(*k)})
}

func (k *Key) Name() string {
	return k.Note.Name + k.Suffix
}

type KeyKind string

const (
	KeyMajor KeyKind = "M"
	KeyMinor KeyKind = "m"
)

func GenerateKeys(kind KeyKind) []Key {
	return std.GenerateKeys(kind)
}

func (t *Theory) GenerateKeys(kind KeyKind) []Key {
	keys := make([]Key, 0, len(t.Config.NaturalNoteNames)*3)

	suffix := ""
	if kind == KeyMinor {
		suffix = string(t.Config.MinorKeySymbols[0])
	}

	for i, nnn := range t.Config.NaturalNoteNames {

		degreeClass := DegreeClass(i)
		pitchClass := t.Config.PitchClassFromDegreeClass(degreeClass)

		keys = append(keys, Key{
			Note: Note{
				Name:        string(nnn),
				DegreeClass: degreeClass,
				PitchClass:  pitchClass,
				Accidentals: 0,
			},
			Suffix: suffix,
			Kind:   kind,
		})

		keys = append(keys, Key{
			Note: Note{
				Name:        string(nnn) + string(t.Config.SharpSymbols[0]),
				DegreeClass: degreeClass,
				PitchClass:  t.Config.AdjustPitchClass(pitchClass, 1),
				Accidentals: 1,
			},
			Suffix: suffix,
			Kind:   kind,
		})

		keys = append(keys, Key{
			Note: Note{
				Name:        string(nnn) + string(t.Config.FlatSymbols[0]),
				DegreeClass: degreeClass,
				PitchClass:  t.Config.AdjustPitchClass(pitchClass, -1),
				Accidentals: -1,
			},
			Suffix: suffix,
			Kind:   kind,
		})
	}

	SortKeys(keys)

	return keys
}

func InferKey(chords []Chord) []Key {
	return std.InferKey(chords)
}

func (t *Theory) InferKey(chords []Chord) []Key {
	panic("not implemented")
}

func MustKey(key Key, err error) Key {
	if err != nil {
		panic(err)
	}

	return key
}

func ParseKey(text string) (Key, error) {
	return std.ParseKey(text)
}

func (t *Theory) ParseKey(text string) (Key, error) {
	n, pos, err := t.parseNote(text, 0)
	if err != nil {
		return Key{}, fmt.Errorf("expected note at position 0: %w", err)
	}

	kind := KeyMajor
	suffix := ""
	if len(text) > pos {
		for _, sym := range t.Config.MinorKeySymbols {
			if strings.HasPrefix(text[pos:], sym) {
				kind = KeyMinor
				suffix = string(sym)
				pos += len(sym)
				break
			}
		}
	}

	if len(text) != pos {
		return Key{}, fmt.Errorf("expected EOF at position %d, but had %q", pos, text[pos:])
	}

	return Key{
		Note:   n,
		Suffix: suffix,
		Kind:   kind,
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
		case keys[i].Suffix < keys[j].Suffix:
			return true
		case keys[i].Suffix > keys[j].Suffix:
			return false
		default:
			return false
		}
	})
}

func TransposeKey(key Key, interval Interval) Key {
	return std.TransposeKey(key, interval)
}

func (t *Theory) TransposeKey(key Key, interval Interval) Key {
	newKeyNote := t.TransposeNote(key.Note, interval)
	return Key{
		Note: newKeyNote,
		Kind: key.Kind,
	}
}
