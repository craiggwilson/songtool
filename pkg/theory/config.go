package theory

import (
	"sort"
)

var defaultConfig Config = DefaultConfig()

func DefaultConfig() Config {
	modGroups := []ChordModifierGroup{
		{
			Name: "Quality",
			Modifiers: []ChordModifier{
				{
					Match:     "m",
					Except:    "maj",
					Semitones: []int{-4, 3},
				},
				{
					Match:     "dim",
					Semitones: []int{-4, 3, -7, 6},
				},
				{
					Match:     "aug",
					Semitones: []int{-7, 8},
				},
			},
		},
		{
			Name: "Numbered",
			Modifiers: []ChordModifier{
				{
					Match:     "maj7",
					Semitones: []int{11},
				},
				{
					Match:     "maj9",
					Semitones: []int{11, 14},
				},
				{
					Match:     "maj11",
					Semitones: []int{11, 14, 17},
				},
				{
					Match:     "maj13",
					Semitones: []int{11, 14, 17, 21},
				},
				{
					Match:     "maj13",
					Semitones: []int{11, 14, 17, 21},
				},
				{
					Match:     "2",
					Semitones: []int{2, -3, -4},
				},
				{
					Match:     "5",
					Semitones: []int{-3, -4},
				},
				{
					Match:     "6",
					Semitones: []int{9},
				},
				{
					Match:     "7",
					Semitones: []int{10},
				},
				{
					Match:     "9",
					Semitones: []int{10, 14},
				},
				{
					Match:     "11",
					Semitones: []int{10, 14, 17},
				},
				{
					Match:     "13",
					Semitones: []int{10, 14, 17, 21},
				},
			},
		},
		{
			Name: "Suspensions",
			Modifiers: []ChordModifier{
				{
					Match:     "sus2",
					Semitones: []int{2, -3, -4},
				},
				{
					Match:     "sus4",
					Semitones: []int{5, -3, -4},
				},
				{
					Match:     "sus",
					Semitones: []int{5, -3, -4},
				},
			},
		},
		{
			Name: "Added Tones",
			Modifiers: []ChordModifier{
				{
					Match:     "add2",
					Semitones: []int{2},
				},
				{
					Match:     "add4",
					Semitones: []int{5},
				},
				{
					Match:     "add6",
					Semitones: []int{9},
				},
				{
					Match:     "add9",
					Semitones: []int{14},
				},
			},
		},
	}

	for _, grp := range modGroups {
		sort.Slice(grp.Modifiers, func(i, j int) bool {
			return len(grp.Modifiers[i].Match) > len(grp.Modifiers[j].Match)
		})
	}

	return Config{
		MinorKeySymbols:         []rune{'m'},
		NaturalNoteNames:        []rune{'C', 'D', 'E', 'F', 'G', 'A', 'B'},
		SharpSymbols:            []rune{'#'},
		FlatSymbols:             []rune{'b'},
		BaseNoteDelimiters:      []rune{'/'},
		MajorChordIntervals:     []int{1, 4, 7},
		ChordModifiers:          modGroups,
		PitchClassCount:         12,
		DegreeClassToPitchClass: []PitchClass{0, 2, 4, 5, 7, 9, 11},
	}
}

type Config struct {
	MinorKeySymbols    []rune
	NaturalNoteNames   []rune
	SharpSymbols       []rune
	FlatSymbols        []rune
	BaseNoteDelimiters []rune

	MajorChordIntervals []int

	ChordModifiers []ChordModifierGroup

	PitchClassCount         int
	DegreeClassToPitchClass []PitchClass
}

type ChordModifierGroup struct {
	Name      string
	Modifiers []ChordModifier
}

type ChordModifier struct {
	Match     string
	Except    string
	Semitones []int
}
