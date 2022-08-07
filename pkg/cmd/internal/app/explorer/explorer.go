package explorer

import tea "github.com/charmbracelet/bubbletea"

type Model struct {
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return "\n\nFile View\n\n"
}
