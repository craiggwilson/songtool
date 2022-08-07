package message

import tea "github.com/charmbracelet/bubbletea"

func ChangeHelpMode(full bool) tea.Cmd {
	return func() tea.Msg { return ChangeHelpModeMsg{Full: full} }
}

type ChangeHelpModeMsg struct {
	Full bool
}
