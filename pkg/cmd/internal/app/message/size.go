package message

import tea "github.com/charmbracelet/bubbletea"

func Invalidate() tea.Cmd {
	return func() tea.Msg {
		return InvalidateMsg{}
	}
}

type InvalidateMsg struct{}
