package models

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/config"
)

var (
	// Text styles.
	chordStyle                = lipgloss.NewStyle()
	directiveStyle            = lipgloss.NewStyle()
	lyricsStyle               = lipgloss.NewStyle()
	sectionNameStyle          = lipgloss.NewStyle()
	headerFooterBoundaryStyle = lipgloss.NewStyle()
	titleStyle                = lipgloss.NewStyle()

	// Layout styles.
	colStyle     = lipgloss.NewStyle().MarginLeft(2)
	contentStyle = lipgloss.NewStyle().Align(lipgloss.Center)

	sectionStyle = lipgloss.NewStyle()

	headerStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	footerStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return headerStyle.Copy().BorderStyle(b)
	}()
)

func initStyles(styles *config.Styles) {
	headerStyle = headerStyle.BorderForeground(styles.BoundaryColor.Color())
	headerFooterBoundaryStyle = headerFooterBoundaryStyle.Foreground(styles.BoundaryColor.Color())
	footerStyle = footerStyle.BorderForeground(styles.BoundaryColor.Color())

	chordStyle = styles.Chord.Apply(chordStyle)
	directiveStyle = styles.Directive.Apply(directiveStyle)
	lyricsStyle = styles.Lyrics.Apply(lyricsStyle)
	sectionNameStyle = styles.SectionName.Apply(sectionNameStyle)
	titleStyle = styles.TitleStyle.Apply(titleStyle)
}
