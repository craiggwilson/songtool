package header

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/message"
	"github.com/craiggwilson/songtool/pkg/songio"
)

var (
	titleBorderStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()
)

func New() Model {
	return Model{}
}

type Model struct {
	BorderColor lipgloss.TerminalColor
	KeyStyle    lipgloss.Style
	TitleStyle  lipgloss.Style
	Width       int

	Meta *songio.Meta
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch tmsg := msg.(type) {
	case message.UpdateSongMsg:
		m.Meta = &tmsg.Meta
	}
	return m, nil
}

func (m Model) View() string {
	title := "<no song>"
	if m.Meta != nil {
		title = m.Meta.Title
		if m.Meta.Key != nil {
			title += fmt.Sprintf(" [%s]", m.KeyStyle.Render(m.Meta.Key.Name))
		}
	}

	title = titleBorderStyle.BorderForeground(m.BorderColor).Render(m.TitleStyle.Render(title))

	line := lipgloss.NewStyle().Foreground(m.BorderColor).Render(strings.Repeat("─", max(0, m.Width-lipgloss.Width(title))))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
