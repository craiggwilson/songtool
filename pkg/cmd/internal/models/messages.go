package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/craiggwilson/songtool/pkg/songio"
	"github.com/craiggwilson/songtool/pkg/theory/interval"
	"github.com/craiggwilson/songtool/pkg/theory/note"
)

func StatusError(err error) tea.Cmd {
	return func() tea.Msg {
		return StatusErrorMsg(err)
	}
}

type StatusErrorMsg error

func Transpose(intval interval.Interval) tea.Cmd {
	return func() tea.Msg {
		return TransposeMsg{Interval: intval}
	}
}

type TransposeMsg struct {
	Interval interval.Interval
}

func UpdateSongFromSource(noteNamer note.Namer, path string, src songio.Reader) tea.Cmd {
	return func() tea.Msg {
		lines, err := songio.ReadAllLines(src)
		if err != nil {
			return StatusErrorMsg(err)
		}

		meta, err := songio.ReadMeta(noteNamer, songio.FromLines(lines), false)
		if err != nil {
			return StatusErrorMsg(err)
		}
		if len(meta.Title) == 0 {
			meta.Title = path
		}

		return UpdateSongMsg{
			Meta:  meta,
			Lines: lines,
		}
	}
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
