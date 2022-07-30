package theory

import (
	"regexp"

	"github.com/craiggwilson/songtool/pkg/theory/interval"
)

func DefaultConfigBase() ConfigBase {
	return ConfigBase{
		NaturalNoteNames:   [7]string{"C", "D", "E", "F", "G", "A", "B"},
		SharpSymbols:       []string{"#"},
		FlatSymbols:        []string{"b"},
		MajorSymbols:       []string{"maj", "M"},
		MinorSymbols:       []string{"m", "-"},
		AugmentedSymbols:   []string{"aug", "+"},
		DiminishedSymobls:  []string{"dim", "°", "o"},
		BaseNoteDelimiters: []string{"/"},
		Scales: map[string][]interval.Interval{
			"Major":     interval.Scales.Ionian,
			"Ionian":    interval.Scales.Ionian,
			"Chromatic": interval.Scales.Chromatic,
		},
	}
}

func DefaultConfig() *Config {
	return NewConfig(DefaultConfigBase())
}

func NewConfig(base ConfigBase) *Config {
	return &Config{
		ConfigBase:     base,
		ChordModifiers: BuildChordModifiers(base),
	}
}

type ConfigBase struct {
	NaturalNoteNames   [7]string                      `json:"naturalNoteNames"`
	SharpSymbols       []string                       `json:"sharpSymbols"`
	FlatSymbols        []string                       `json:"flatSymbols"`
	MajorSymbols       []string                       `json:"majorSymbols"`
	MinorSymbols       []string                       `json:"minorSymbols"`
	AugmentedSymbols   []string                       `json:"augmentedSymbols"`
	DiminishedSymobls  []string                       `json:"diminishedSymbols"`
	BaseNoteDelimiters []string                       `json:"baseNoteDelimiters"`
	Scales             map[string][]interval.Interval `json:"scales"`
}

type Config struct {
	ConfigBase

	ChordModifiers []ChordModifier `json:"chordModifiers"`
}

type ChordModifier struct {
	Name   string              `json:"name"`
	Match  *regexp.Regexp      `json:"match"`
	Except *regexp.Regexp      `json:"except"`
	Add    []interval.Interval `json:"add"`
	Remove []interval.Interval `json:"remove"`
}

func BuildChordModifiers(cfg ConfigBase) []ChordModifier {
	return []ChordModifier{
		{
			Name: "Base",
			Add:  []interval.Interval{interval.Perfect(0), interval.Major(2), interval.Perfect(4)},
		},
		{
			Name:   "Minor",
			Match:  regexp.MustCompile("^" + regexOr(cfg.MinorSymbols)),
			Except: regexp.MustCompile("^" + regexOr(cfg.MajorSymbols)),
			Add:    []interval.Interval{interval.Minor(2)},
			Remove: []interval.Interval{interval.Major(2)},
		},
		{
			Name:   "Augmented",
			Match:  regexp.MustCompile("^" + regexOr(cfg.AugmentedSymbols)),
			Add:    []interval.Interval{interval.Augmented(4, 1)},
			Remove: []interval.Interval{interval.Perfect(4)},
		},
		{
			Name:   "Diminished",
			Match:  regexp.MustCompile("^" + regexOr(cfg.DiminishedSymobls)),
			Add:    []interval.Interval{interval.Minor(2), interval.Diminished(4, 1)},
			Remove: []interval.Interval{interval.Major(2), interval.Perfect(4)},
		},
		{
			Name:   "2nd (alt for sus2)",
			Match:  regexp.MustCompile("^(" + regexOr(cfg.MinorSymbols, cfg.DiminishedSymobls, cfg.AugmentedSymbols) + ")?(?P<mod>2)"),
			Add:    []interval.Interval{interval.Major(1)},
			Remove: []interval.Interval{interval.Minor(2), interval.Major(2)},
		},
		{
			Name:   "5th (no 3rd)",
			Match:  regexp.MustCompile("^5"),
			Remove: []interval.Interval{interval.Major(2)},
		},
		{
			Name:   "6th",
			Match:  regexp.MustCompile("^(" + regexOr(cfg.MinorSymbols, cfg.DiminishedSymobls, cfg.AugmentedSymbols) + ")?(?P<mod>6)"),
			Except: regexp.MustCompile("^(" + regexOr(cfg.MinorSymbols, cfg.DiminishedSymobls, cfg.AugmentedSymbols) + ")?(?P<mod>69)"),
			Add:    []interval.Interval{interval.Major(5)},
		},
		{
			Name:  "6th+9th",
			Match: regexp.MustCompile("^(" + regexOr(cfg.MinorSymbols, cfg.DiminishedSymobls, cfg.AugmentedSymbols) + ")?(?P<mod>69)"),
			Add:   []interval.Interval{interval.Major(5), interval.Major(8)},
		},
		{
			Name:  "7th",
			Match: regexp.MustCompile("^(" + regexOr(cfg.MinorSymbols, cfg.AugmentedSymbols) + ")?(?P<mod>7)"),
			Add:   []interval.Interval{interval.Minor(6)},
		},
		{
			Name:  "Diminished 7th",
			Match: regexp.MustCompile("^(" + regexOr(cfg.DiminishedSymobls) + ")(?P<mod>7)"),
			Add:   []interval.Interval{interval.Diminished(6, 1)},
		},
		{
			Name:  "9th",
			Match: regexp.MustCompile("^(" + regexOr(cfg.MinorSymbols, cfg.DiminishedSymobls, cfg.AugmentedSymbols) + ")?(?P<mod>9)"),
			Add:   []interval.Interval{interval.Minor(6), interval.Major(8)},
		},
		{
			Name:  "11th",
			Match: regexp.MustCompile("^(" + regexOr(cfg.MinorSymbols, cfg.DiminishedSymobls, cfg.AugmentedSymbols) + ")?(?P<mod>11)"),
			Add:   []interval.Interval{interval.Minor(6), interval.Major(8), interval.Perfect(10)},
		},
		{
			Name:  "13th",
			Match: regexp.MustCompile("^(" + regexOr(cfg.MinorSymbols, cfg.DiminishedSymobls, cfg.AugmentedSymbols) + ")?(?P<mod>13)"),
			Add:   []interval.Interval{interval.Minor(6), interval.Major(8), interval.Perfect(10), interval.Major(12)},
		},
		{
			Name:   "Major 7th",
			Match:  regexp.MustCompile(regexOrWithSuffix("7", cfg.MajorSymbols)),
			Add:    []interval.Interval{interval.Major(6)},
			Remove: []interval.Interval{interval.Minor(6)},
		},
		{
			Name:   "Major 9th",
			Match:  regexp.MustCompile(regexOrWithSuffix("9", cfg.MajorSymbols)),
			Add:    []interval.Interval{interval.Major(6), interval.Major(8)},
			Remove: []interval.Interval{interval.Minor(6)},
		},
		{
			Name:   "Major 11th",
			Match:  regexp.MustCompile(regexOrWithSuffix("11", cfg.MajorSymbols)),
			Add:    []interval.Interval{interval.Major(6), interval.Major(8), interval.Perfect(10)},
			Remove: []interval.Interval{interval.Minor(6)},
		},
		{
			Name:   "Major 13th",
			Match:  regexp.MustCompile(regexOrWithSuffix("13", cfg.MajorSymbols)),
			Add:    []interval.Interval{interval.Major(6), interval.Major(8), interval.Perfect(10), interval.Major(12)},
			Remove: []interval.Interval{interval.Minor(6)},
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
			Match:  regexp.MustCompile(`\(` + regexOrWithSuffix("5", cfg.FlatSymbols) + `\)|` + regexOrWithSuffix("5", cfg.FlatSymbols)),
			Add:    []interval.Interval{interval.Diminished(4, 1)},
			Remove: []interval.Interval{interval.Perfect(4)},
		},
		{
			Name:  "Flat 6th",
			Match: regexp.MustCompile(`\(` + regexOrWithSuffix("6", cfg.FlatSymbols) + `\)|` + regexOrWithSuffix("6", cfg.FlatSymbols)),
			Add:   []interval.Interval{interval.Minor(5)},
		},
		{
			Name:  "Flat 9th",
			Match: regexp.MustCompile(`\(` + regexOrWithSuffix("9", cfg.FlatSymbols) + `\)|` + regexOrWithSuffix("9", cfg.FlatSymbols)),
			Add:   []interval.Interval{interval.Minor(8)},
		},
		{
			Name:  "Flat 13th",
			Match: regexp.MustCompile(`\(` + regexOrWithSuffix("13", cfg.FlatSymbols) + `\)|` + regexOrWithSuffix("13", cfg.FlatSymbols)),
			Add:   []interval.Interval{interval.Minor(8)},
		},
		{
			Name:   "Sharp 4th",
			Match:  regexp.MustCompile(`\(` + regexOrWithSuffix("4", cfg.SharpSymbols) + `\)|` + regexOrWithSuffix("4", cfg.SharpSymbols)),
			Add:    []interval.Interval{interval.Augmented(3, 1)},
			Remove: []interval.Interval{interval.Perfect(4)},
		},
		{
			Name:   "Sharp 5th",
			Match:  regexp.MustCompile(`\(` + regexOrWithSuffix("5", cfg.SharpSymbols) + `\)|` + regexOrWithSuffix("5", cfg.SharpSymbols)),
			Add:    []interval.Interval{interval.Augmented(4, 1)},
			Remove: []interval.Interval{interval.Perfect(4)},
		},
		{
			Name:  "Sharp 9th",
			Match: regexp.MustCompile(`\(` + regexOrWithSuffix("9", cfg.SharpSymbols) + `\)|` + regexOrWithSuffix("9", cfg.SharpSymbols)),
			Add:   []interval.Interval{interval.Augmented(8, 1)},
		},
		{
			Name:  "Sharp 11th",
			Match: regexp.MustCompile(`\(` + regexOrWithSuffix("11", cfg.SharpSymbols) + `\)|` + regexOrWithSuffix("11", cfg.SharpSymbols)),
			Add:   []interval.Interval{interval.Augmented(10, 1)},
		},
	}
}

func regexOr(sss ...[]string) string {
	result := ""
	count := 0
	for _, ss := range sss {
		for _, s := range ss {
			if len(s) > 0 {
				if count > 0 {
					result += "|"
				}
				result += regexp.QuoteMeta(s)
				count++
			}
		}
	}

	return result
}

func regexOrWithSuffix(suffix string, sss ...[]string) string {
	result := ""
	count := 0
	for _, ss := range sss {
		for _, s := range ss {
			if len(s) > 0 {
				if count > 0 {
					result += "|"
				}
				result += regexp.QuoteMeta(s) + suffix
				count++
			}
		}
	}

	return result
}
