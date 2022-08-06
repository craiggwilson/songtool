package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/craiggwilson/songtool/pkg/songio"
)

func newHeaderViewModel() headerViewModel {
	return headerViewModel{}
}

type headerViewModel struct {
	Width int
	Meta  *songio.Meta
}

func (m headerViewModel) Update(msg tea.Msg) (headerViewModel, tea.Cmd) {
	return m, nil
}

func (m headerViewModel) View() string {
	if m.Meta == nil {
		return "<no song>"
	}

	var title string
	header := m.Meta.Title
	if m.Meta.Key != nil {
		header += fmt.Sprintf(" [%s]", chordStyle.Render(m.Meta.Key.Name))
	}
	title = headerStyle.Render(titleStyle.Render(header))

	line := headerFooterBoundaryStyle.Render(strings.Repeat("─", max(0, m.Width-lipgloss.Width(title))))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}
