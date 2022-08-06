package message

import (
	tea "github.com/charmbracelet/bubbletea"
)

func UpdateStatusInfo(text string) tea.Cmd {
	return func() tea.Msg {
		return UpdateStatusMsg{Info: text}
	}
}

func UpdateStatusError(err error) tea.Cmd {
	return func() tea.Msg {
		return UpdateStatusMsg{Err: err}
	}
}

type UpdateStatusMsg struct {
	Info string
	Err  error
}
