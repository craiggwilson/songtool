package internal

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

type color bool

func (c color) AfterApply() error {
	if !c {
		lipgloss.SetColorProfile(termenv.Ascii)
	} else if termenv.EnvColorProfile() == termenv.Ascii {
		lipgloss.SetColorProfile(termenv.ANSI)
	}

	return nil
}
