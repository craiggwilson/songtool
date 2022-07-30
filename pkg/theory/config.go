package theory

import (
	"regexp"

	"github.com/craiggwilson/songtool/pkg/theory/interval"
)

const (
	patternMinor = "m"
	patternMajor = "maj"
	patternDim   = "dim"
	patternAug   = "aug"
	numPrefix    = "m(aj)?|dim|aug"
)

func DefaultConfig() *Config {
	return &Config{
		NaturalNoteNames:   [7]string{"C", "D", "E", "F", "G", "A", "B"},
		SharpSymbols:       []string{"#"},
		FlatSymbols:        []string{"b"},
		MinorKeySymbols:    []string{"m"},
		BaseNoteDelimiters: []string{"/"},
		Scales: map[string][]interval.Interval{
			"Major":     interval.Scales.Ionian,
			"Ionian":    interval.Scales.Ionian,
			"Chromatic": interval.Scales.Chromatic,
		},
		ChordModifiers: []ChordModifier{
			{
				Name: "Base",
				Add:  []interval.Interval{interval.Perfect(0), interval.Major(2), interval.Perfect(4)},
			},
			{
				Name:   "Minor",
				Match:  regexp.MustCompile("^m"),
				Except: regexp.MustCompile("^maj"),
				Add:    []interval.Interval{interval.Minor(2)},
				Remove: []interval.Interval{interval.Major(2)},
			},
			{
				Name:   "Augmented",
				Match:  regexp.MustCompile("^aug"),
				Add:    []interval.Interval{interval.Augmented(4, 1)},
				Remove: []interval.Interval{interval.Perfect(4)},
			},
			{
				Name:   "Diminished",
				Match:  regexp.MustCompile("^dim"),
				Add:    []interval.Interval{interval.Minor(2), interval.Diminished(4, 1)},
				Remove: []interval.Interval{interval.Major(2), interval.Perfect(4)},
			},
			{
				Name:   "2nd (alt for sus2)",
				Match:  regexp.MustCompile(`^(m|dim|aug)?(?P<mod>2)`),
				Add:    []interval.Interval{interval.Major(1)},
				Remove: []interval.Interval{interval.Minor(2), interval.Major(2)},
			},
			// {
			// 	Name:   "4th (alt for sus)",
			// 	Match:  regexp.MustCompile("^(m|dim|aug)?(?P<mod>4)"),
			// 	Add:    []interval.Interval{interval.Major(1)},
			// 	Remove: []interval.Interval{interval.Major(2)},
			// },
			{
				Name:   "5th (no 3rd)",
				Match:  regexp.MustCompile("^5"),
				Remove: []interval.Interval{interval.Major(2)},
			},
			{
				Name:   "6th",
				Match:  regexp.MustCompile("^(m|dim|aug)?(?P<mod>6)"),
				Except: regexp.MustCompile("^(m|dim|aug)?69"),
				Add:    []interval.Interval{interval.Major(5)},
			},
			{
				Name:  "6th+9th",
				Match: regexp.MustCompile("^(m|dim|aug)?(?P<mod>69)"),
				Add:   []interval.Interval{interval.Major(5), interval.Major(8)},
			},
			{
				Name:  "7th",
				Match: regexp.MustCompile("^(m|aug)?(?P<mod>7)"),
				Add:   []interval.Interval{interval.Minor(6)},
			},
			{
				Name:  "Diminished 7th",
				Match: regexp.MustCompile("^dim(7)"),
				Add:   []interval.Interval{interval.Diminished(6, 1)},
			},
			{
				Name:  "9th",
				Match: regexp.MustCompile("^(m|dim|aug)?(?P<mod>9)"),
				Add:   []interval.Interval{interval.Minor(6), interval.Major(8)},
			},
			{
				Name:  "11th",
				Match: regexp.MustCompile("^(m|dim|aug)?(?P<mod>11)"),
				Add:   []interval.Interval{interval.Minor(6), interval.Major(8), interval.Perfect(10)},
			},
			{
				Name:  "13th",
				Match: regexp.MustCompile("^(m|dim|aug)?(?P<mod>13)"),
				Add:   []interval.Interval{interval.Minor(6), interval.Major(8), interval.Perfect(10), interval.Major(12)},
			},
			{
				Name:  "Major 7th",
				Match: regexp.MustCompile("maj7"),
				Add:   []interval.Interval{interval.Major(6)},
			},
			{
				Name:  "Major 9th",
				Match: regexp.MustCompile("maj9"),
				Add:   []interval.Interval{interval.Major(6), interval.Major(8)},
			},
			{
				Name:  "Major 11th",
				Match: regexp.MustCompile("maj11"),
				Add:   []interval.Interval{interval.Major(6), interval.Major(8), interval.Perfect(10)},
			},
			{
				Name:  "Major 13th",
				Match: regexp.MustCompile("maj13"),
				Add:   []interval.Interval{interval.Major(6), interval.Major(8), interval.Perfect(10), interval.Major(12)},
			},
			{
				Name:   "Suspended 2nd",
				Match:  regexp.MustCompile("sus2"),
				Add:    []interval.Interval{interval.Major(1)},
				Remove: []interval.Interval{interval.Major(2)},
			},
			{
				Name:   "Suspended 4th",
				Match:  regexp.MustCompile("sus4?"),
				Except: regexp.MustCompile("sus2"),
				Add:    []interval.Interval{interval.Perfect(3)},
				Remove: []interval.Interval{interval.Minor(2), interval.Major(2)},
			},
			{
				Name:  "Added 2nd",
				Match: regexp.MustCompile("add2"),
				Add:   []interval.Interval{interval.Major(1)},
			},
			{
				Name:  "Added 4th",
				Match: regexp.MustCompile("add4"),
				Add:   []interval.Interval{interval.Perfect(3)},
			},
			{
				Name:  "Added 6th",
				Match: regexp.MustCompile("add6"),
				Add:   []interval.Interval{interval.Major(5)},
			},
			{
				Name:  "Added 9th",
				Match: regexp.MustCompile("add9|/9"),
				Add:   []interval.Interval{interval.Major(8)},
			},
			{
				Name:  "Added 11th",
				Match: regexp.MustCompile("add11"),
				Add:   []interval.Interval{interval.Perfect(10)},
			},
			{
				Name:  "Added 13th",
				Match: regexp.MustCompile("add13"),
				Add:   []interval.Interval{interval.Major(12)},
			},
			{
				Name:   "Flat 5th",
				Match:  regexp.MustCompile(`\(b5\)|b5`),
				Add:    []interval.Interval{interval.Diminished(4, 1)},
				Remove: []interval.Interval{interval.Perfect(4)},
			},
			{
				Name:  "Flat 6th",
				Match: regexp.MustCompile(`\(b6\)|\(b13\)|b6|b13`),
				Add:   []interval.Interval{interval.Minor(5)},
			},
			{
				Name:  "Flat 9th",
				Match: regexp.MustCompile(`\(b9\)|b9`),
				Add:   []interval.Interval{interval.Minor(1)},
			},
			{
				Name:  "Sharp 5th",
				Match: regexp.MustCompile(`\(#5\)|#5`),
				Add:   []interval.Interval{interval.Augmented(4, 1)},
			},
			{
				Name:  "Sharp 9th",
				Match: regexp.MustCompile(`\(#9\)|#9`),
				Add:   []interval.Interval{interval.Augmented(2, 1)},
			},
			{
				Name:  "Sharp 11th",
				Match: regexp.MustCompile(`\(#4\)|\(#11\)|#4|#11`),
				Add:   []interval.Interval{interval.Augmented(3, 1)},
			},
		},
	}
}

type Config struct {
	NaturalNoteNames   [7]string                      `json:"naturalNoteNames"`
	SharpSymbols       []string                       `json:"sharpSymbols"`
	FlatSymbols        []string                       `json:"flatSymbols"`
	MajorKeySymbols    []string                       `json:"majorKeySymbols"`
	MinorKeySymbols    []string                       `json:"minorKeySymbols"`
	BaseNoteDelimiters []string                       `json:"baseNoteDelimiters"`
	Scales             map[string][]interval.Interval `json:"scales"`

	ChordModifiers []ChordModifier `json:"chordMofifiers"`
}

type ChordModifier struct {
	Name   string              `json:"name"`
	Match  *regexp.Regexp      `json:"match"`
	Except *regexp.Regexp      `json:"except"`
	Add    []interval.Interval `json:"add"`
	Remove []interval.Interval `json:"remove"`
}
