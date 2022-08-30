package app

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

type mode int

func (m mode) IsCommandMode() bool {
	return m&modeCommand == modeCommand
}

func (m mode) IsExplorerMode() bool {
	return m&modeExplorer == modeExplorer
}

func (m mode) IsSongMode() bool {
	return m&modeSong == modeSong
}

const (
	modeSong            = 1
	modeSongCommand     = modeSong | modeCommand
	modeExplorer        = 2
	modeExplorerCommand = modeExplorer | modeCommand

	modeCommand = 4
)
