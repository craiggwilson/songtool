package cmd

import (
	"github.com/craiggwilson/songtool/pkg/theory"
	"github.com/fatih/color"
)

func LoadConfig(path string) *Config {
	return &Config{
		Highlighter: Highlighter{
			DirectiveColor: color.New(color.FgBlue),
			SectionColor:   color.New(color.FgCyan),
			ChordColor:     color.New(color.FgHiCyan),
		},
		Theory: theory.DefaultConfig(),
	}
}

type Config struct {
	Highlighter Highlighter   `json:"highlighter"`
	Theory      theory.Config `json:"theory"`
}

type Highlighter struct {
	DirectiveColor *color.Color `json:"directive"`
	SectionColor   *color.Color `json:"section"`
	ChordColor     *color.Color `json:"chord"`
}
