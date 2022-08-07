package message

import tea "github.com/charmbracelet/bubbletea"

func EnterCommandMode(value string) tea.Cmd {
	return func() tea.Msg { return EnterCommandModeMsg{value} }
}

type EnterCommandModeMsg struct {
	Value string
}

func ExitCommandMode() tea.Cmd {
	return func() tea.Msg { return ExitCommandModeMsg{} }
}

type ExitCommandModeMsg struct{}

func EnterSongMode() tea.Cmd {
	return func() tea.Msg { return EnterSongModeMsg{} }
}

type EnterSongModeMsg struct{}

func EnterExplorerMode() tea.Cmd {
	return func() tea.Msg { return EnterExplorerModeMsg{} }
}

type EnterExplorerModeMsg struct{}
