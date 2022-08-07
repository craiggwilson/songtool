package message

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/craiggwilson/songtool/pkg/songio"
	"github.com/craiggwilson/songtool/pkg/theory/interval"
)

func LoadSong(path string) tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			return OpenSongMsg{Path: path}
		},
		EnterSongMode(),
	)
}

type OpenSongMsg struct {
	Path string
}

func TransposeSong(intval interval.Interval) tea.Cmd {
	return func() tea.Msg {
		return TransposeSongMsg{Interval: intval}
	}
}

type TransposeSongMsg struct {
	Interval interval.Interval
}

func UpdateSong(meta songio.Meta, lines []songio.Line) tea.Cmd {
	return func() tea.Msg {
		return UpdateSongMsg{
			Meta:  meta,
			Lines: lines,
		}
	}
}

type UpdateSongMsg struct {
	Meta  songio.Meta
	Lines []songio.Line
}
