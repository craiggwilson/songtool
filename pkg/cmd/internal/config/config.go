package config

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/craiggwilson/songtool/pkg/theory"
)

type Config struct {
	File
	Theory *theory.Theory

	Path string
}

type File struct {
	Edit   Edit              `json:"edit"`
	Styles Styles            `json:"styles"`
	Theory theory.ConfigBase `json:"theory"`
}

type Edit struct {
	Command string   `json:"command"`
	Args    []string `json:"args,omitempty"`
}

type Styles struct {
	MaxColumns int `json:"maxColumns,omitempty"`

	BoundaryColor Color `json:"boundaryColor,omitempty"`
	Chord         Style `json:"chord,omitempty"`
	Directive     Style `json:"directive,omitempty"`
	Lyrics        Style `json:"lyrics,omitempty"`
	SectionName   Style `json:"sectionName,omitempty"`
	TitleStyle    Style `json:"titleStyle,omitempty"`
}

type Style struct {
	Background Color `json:"background,omitempty"`
	Bold       bool  `json:"bold"`
	Foreground Color `json:"foreground,omitempty"`
	Italic     bool  `json:"italic"`
	Underline  bool  `json:"underline"`

	style *lipgloss.Style
}

type Color struct {
	Light string `json:"light,omitempty"`
	Dark  string `json:"dark,omitempty"`
}

func (csc Color) Color() lipgloss.TerminalColor {
	switch {
	case len(csc.Light) != 0 && len(csc.Dark) != 0:
		return lipgloss.AdaptiveColor{
			Light: csc.Light,
			Dark:  csc.Dark,
		}
	case len(csc.Light) != 0:
		return lipgloss.Color(csc.Light)
	case len(csc.Dark) != 0:
		return lipgloss.Color(csc.Dark)
	default:
		return lipgloss.NoColor{}
	}
}

func (s *Style) Apply(style lipgloss.Style) lipgloss.Style {
	return style.Background(s.Background.Color()).
		Bold(s.Bold).
		Foreground(s.Foreground.Color()).
		Italic(s.Italic).
		Underline(s.Underline)
}

func (s *Style) Render(text string) string {
	if s.style == nil {
		ns := s.Apply(lipgloss.NewStyle())
		s.style = &ns
	}

	return s.style.Render(text)
}
