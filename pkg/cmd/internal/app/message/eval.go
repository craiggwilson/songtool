package message

import (
	tea "github.com/charmbracelet/bubbletea"
)

func Eval(text string) tea.Cmd {
	return func() tea.Msg {
		return EvalMsg{
			Text: text,
		}
	}
}

type EvalMsg struct {
	Text string
}
