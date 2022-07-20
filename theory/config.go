package theory

import (
	"sort"

	"github.com/craiggwilson/songtools/theory/note"
)

const DegreeClassCount = 7
const PitchClassCount = 12

var degreeClassToPitchClass = []note.PitchClass{0, 2, 4, 5, 7, 9, 11}

func DefaultConfig() Config {

	modGroups := []ChordModifierGroup{
		{
			Name: "Quality",
			Modifiers: []ChordModifier{
				{
					Match:     "m",
					Except:    "maj",
					Intervals: []int{-4, 3},
				},
				{
					Match:     "dim",
					Intervals: []int{-4, 3, -7, 6},
				},
				{
					Match:     "aug",
					Intervals: []int{-7, 8},
				},
			},
		},
		{
			Name: "Numbered",
			Modifiers: []ChordModifier{
				{
					Match:     "maj7",
					Intervals: []int{11},
				},
				{
					Match:     "maj9",
					Intervals: []int{11, 14},
				},
				{
					Match:     "maj11",
					Intervals: []int{11, 14, 17},
				},
				{
					Match:     "maj13",
					Intervals: []int{11, 14, 17, 21},
				},
				{
					Match:     "maj13",
					Intervals: []int{11, 14, 17, 21},
				},
				{
					Match:     "2",
					Intervals: []int{2, -3, -4},
				},
				{
					Match:     "5",
					Intervals: []int{-3, -4},
				},
				{
					Match:     "6",
					Intervals: []int{9},
				},
				{
					Match:     "7",
					Intervals: []int{10},
				},
				{
					Match:     "9",
					Intervals: []int{10, 14},
				},
				{
					Match:     "11",
					Intervals: []int{10, 14, 17},
				},
				{
					Match:     "13",
					Intervals: []int{10, 14, 17, 21},
				},
			},
		},
		{
			Name: "Suspensions",
			Modifiers: []ChordModifier{
				{
					Match:     "sus2",
					Intervals: []int{2, -3, -4},
				},
				{
					Match:     "sus4",
					Intervals: []int{5, -3, -4},
				},
				{
					Match:     "sus",
					Intervals: []int{5, -3, -4},
				},
			},
		},
		{
			Name: "Added Tones",
			Modifiers: []ChordModifier{
				{
					Match:     "add2",
					Intervals: []int{2},
				},
				{
					Match:     "add4",
					Intervals: []int{5},
				},
				{
					Match:     "add6",
					Intervals: []int{9},
				},
				{
					Match:     "add9",
					Intervals: []int{14},
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
		MinorKeySymbols:     []rune{'m'},
		NaturalNoteNames:    []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G'},
		SharpSymbols:        []rune{'#'},
		FlatSymbols:         []rune{'b'},
		BaseNoteDelimiters:  []rune{'/'},
		MajorChordIntervals: []int{1, 4, 7},
		ChordModifiers:      modGroups,
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
}

func (c *Config) DegreeClass(naturalNoteName rune) (note.DegreeClass, bool) {
	for i, nn := range c.NaturalNoteNames {
		if nn == naturalNoteName {
			return note.DegreeClass(i), true
		}
	}

	return 0, false
}

func (C *Config) PitchClassFromDegreeClass(degreeClass note.DegreeClass) (note.PitchClass, bool) {
	if int(degreeClass) < len(degreeClassToPitchClass) {
		return degreeClassToPitchClass[int(degreeClass)], true
	}

	return 0, false
}

type ChordModifierGroup struct {
	Name      string
	Modifiers []ChordModifier
}

type ChordModifier struct {
	Match     string
	Except    string
	Intervals []int
}
