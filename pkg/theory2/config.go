package theory2

import (
	"github.com/craiggwilson/songtool/pkg/theory2/interval"
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
	}
}

type Config struct {
	NaturalNoteNames   [7]string `json:"naturalNoteNames"`
	SharpSymbols       []string  `json:"sharpSymbols"`
	FlatSymbols        []string  `json:"flatSymbols"`
	MajorKeySymbols    []string  `json:"majorKeySymbols"`
	MinorKeySymbols    []string  `json:"minorKeySymbols"`
	BaseNoteDelimiters []string  `json:"baseNoteDelimiters"`

	Scales map[string][]interval.Interval `json:"scales"`
}
