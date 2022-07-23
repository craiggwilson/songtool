package theory

import (
	"fmt"
	"math"
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
		MinorKeySymbols:         []string{"m"},
		NaturalNoteNames:        []string{"C", "D", "E", "F", "G", "A", "B"},
		SharpSymbols:            []string{"#"},
		FlatSymbols:             []string{"b"},
		BaseNoteDelimiters:      []string{"/"},
		MajorChordIntervals:     []int{1, 4, 7},
		ChordModifiers:          modGroups,
		PitchClassCount:         12,
		DegreeClassToPitchClass: []PitchClass{0, 2, 4, 5, 7, 9, 11},
	}
}

type Config struct {
	MinorKeySymbols    []string
	NaturalNoteNames   []string
	SharpSymbols       []string
	FlatSymbols        []string
	BaseNoteDelimiters []string

	MajorChordIntervals []int

	ChordModifiers []ChordModifierGroup

	PitchClassCount         int
	DegreeClassToPitchClass []PitchClass
}

func (cfg *Config) AdjustDegreeClass(degreeClass DegreeClass, by int) DegreeClass {
	return (degreeClass + DegreeClass(by) + DegreeClass(len(cfg.NaturalNoteNames))) % DegreeClass(len(cfg.NaturalNoteNames))
}

func (cfg *Config) AdjustPitchClass(pitchClass PitchClass, by int) PitchClass {
	return (pitchClass + PitchClass(by) + PitchClass(cfg.PitchClassCount)) % PitchClass(cfg.PitchClassCount)
}

func (cfg *Config) DegreeClassDelta(a, b DegreeClass) int {
	return classDelta(int(a), int(b), len(cfg.NaturalNoteNames))
}

func (cfg *Config) DegreeClassFromNaturalNoteName(naturalNoteName string) DegreeClass {
	for i, nn := range cfg.NaturalNoteNames {
		if nn == naturalNoteName {
			return DegreeClass(i)
		}
	}

	panic(fmt.Sprintf("natural note name %q does not map to a degree class", naturalNoteName))
}

func (cfg *Config) DegreeClassFromPitchClass(pitchClass PitchClass, enharmonic Enharmonic) DegreeClass {
	switch enharmonic {
	case Sharp:
		for i := len(cfg.DegreeClassToPitchClass) - 1; i >= 0; i-- {
			if pitchClass >= cfg.DegreeClassToPitchClass[i] {
				return DegreeClass(i)
			}
		}
	case Flat:
		for i := 0; i < len(cfg.DegreeClassToPitchClass); i++ {
			if pitchClass <= cfg.DegreeClassToPitchClass[i] {
				return DegreeClass(i)
			}
		}
	default:
		panic(fmt.Sprintf("invalid enharmonic %s", enharmonic))
	}

	panic(fmt.Sprintf("invalid pitch class %d", pitchClass))
}

func (cfg *Config) NormalizeAccidentals(accidentals int) int {
	return normalize(accidentals, cfg.PitchClassCount)
}

func (cfg *Config) PitchClassFromDegreeClass(degreeClass DegreeClass) PitchClass {
	if int(degreeClass) < len(cfg.DegreeClassToPitchClass) {
		return cfg.DegreeClassToPitchClass[int(degreeClass)]
	}

	panic(fmt.Sprintf("degree class %d does not map to a pitch class", degreeClass))
}

func (cfg *Config) PitchClassDelta(a, b PitchClass) int {
	return classDelta(int(a), int(b), cfg.PitchClassCount)
}

func classDelta(a, b, count int) int {
	d1 := b - a
	d2 := b - a + count
	if math.Abs(float64(d1)) < math.Abs(float64(d2)) {
		return d1
	}

	return d2
}

func normalize(v int, count int) int {
	switch {
	case v > count/2:
		return -count + v
	case v < -count/2:
		return count + v
	default:
		return v
	}
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
