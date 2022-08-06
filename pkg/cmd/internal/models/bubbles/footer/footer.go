package footer

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	footerStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()
)

func New() Model {
	return Model{}
}

type Model struct {
	BorderColor        lipgloss.TerminalColor
	ScrollPercentStyle lipgloss.Style
	Width              int

	ScrollPercent float64
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	info := footerStyle.Render(m.ScrollPercentStyle.Render(fmt.Sprintf("%3.f%%", m.ScrollPercent*100)))
	line := lipgloss.NewStyle().Foreground(m.BorderColor).Render(strings.Repeat("─", max(0, m.Width-lipgloss.Width(info))))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
