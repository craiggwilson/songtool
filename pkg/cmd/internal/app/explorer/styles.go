package explorer

import "github.com/charmbracelet/lipgloss"

func DefaultStyles() Styles {
	return Styles{
		ItemStyle:         lipgloss.NewStyle(),
		SelectedItemStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		ColumnStyle:       lipgloss.NewStyle().MarginRight(2),
	}
}

type Styles struct {
	ItemStyle         lipgloss.Style
	SelectedItemStyle lipgloss.Style
	ColumnStyle       lipgloss.Style
}
