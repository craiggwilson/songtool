package app

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/message"
	"github.com/craiggwilson/songtool/pkg/songio"
)

type SongOpenedMsg struct {
	Path  string
	Meta  songio.Meta
	Lines []songio.Line
}

func OpenSong(path string) tea.Cmd {
	return func() tea.Msg {
		return OpenSongMsg{Path: path}
	}
}

type OpenSongMsg struct {
	Path string
}

func SongOpened(meta songio.Meta, lines []songio.Line) tea.Cmd {
	return func() tea.Msg {
		return SongOpenedMsg{
			Meta:  meta,
			Lines: lines,
		}
	}
}

func (m appModel) openSong(path string) tea.Cmd {
	return func() tea.Msg {
		var f *os.File
		var err error
		switch path {
		case "":
			return message.UpdateStatusError(fmt.Errorf("no file to load"))()
		default:
			f, err = os.Open(path)
			if err != nil {
				return message.UpdateStatusError(err)()
			}
		}
		defer f.Close()

		rdr := songio.ReadChordsOverLyrics(m.cfg.Theory, m.cfg.Theory, f)
		lines, err := songio.ReadAllLines(rdr)
		if err != nil {
			return message.UpdateStatusError(err)()
		}

		meta, err := songio.ReadMeta(m.cfg.Theory, songio.FromLines(lines), true)
		if err != nil {
			return message.UpdateStatusError(err)()
		}

		if meta.Title == "" {
			meta.Title = path
		}

		return SongOpened(meta, lines)()
	}
}

type loadedSong struct {
	path string
	meta *songio.Meta

	lines    []songio.Line
	sections []*loadedSongSection
}

func (sm *loadedSong) Sections() []*loadedSongSection {
	return sm.sections
}

type loadedSongSection struct {
}
