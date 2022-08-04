package cmd

import "github.com/charmbracelet/lipgloss"

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
