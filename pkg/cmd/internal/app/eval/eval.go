package eval

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/message"
	"github.com/craiggwilson/songtool/pkg/songio"
	"github.com/craiggwilson/songtool/pkg/theory"
	"github.com/craiggwilson/songtool/pkg/theory/interval"
)

func New(theory *theory.Theory) Model {
	return Model{
		Context: Context{
			Theory: theory,
		},
	}
}

type Model struct {
	Context Context
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch tmsg := msg.(type) {
	case message.EvalMsg:
		return m, run(m.Context, tmsg.Text)
	case message.LoadDirectoryMsg:
		return m, m.listFiles(tmsg.Path)
	case message.OpenSongMsg:
		return m, m.openSong(tmsg.Path)
	case message.TransposeSongMsg:
		return m, m.transposeSong(tmsg.Interval)
	case message.UpdateSongMsg:
		m.Context.Meta = &tmsg.Meta
		m.Context.Lines = tmsg.Lines
	}

	return m, nil
}

func (m Model) listFiles(path string) tea.Cmd {
	return func() tea.Msg {
		entries, err := os.ReadDir(path)
		if err != nil {
			return message.UpdateStatusError(err)()
		}

		files := make([]message.FileItem, len(entries))
		for i := range entries {
			files[i].Path = filepath.Join(path, entries[i].Name())
		}

		return message.UpdateFiles(files)()
	}
}

func (m Model) openSong(path string) tea.Cmd {
	var f *os.File
	var err error
	switch path {
	case "":
		return message.UpdateStatusError(fmt.Errorf("no file to load"))
	default:
		f, err = os.Open(path)
		if err != nil {
			return message.UpdateStatusError(err)
		}
	}
	defer f.Close()

	rdr := songio.ReadChordsOverLyrics(m.Context.Theory, m.Context.Theory, f)
	lines, err := songio.ReadAllLines(rdr)
	if err != nil {
		return message.UpdateStatusError(err)
	}

	meta, err := songio.ReadMeta(m.Context.Theory, songio.FromLines(lines), true)
	if err != nil {
		return message.UpdateStatusError(err)
	}

	if meta.Title == "" {
		meta.Title = path
	}

	return message.UpdateSong(meta, lines)
}

func (m Model) transposeSong(by interval.Interval) tea.Cmd {
	transposed := songio.Transpose(m.Context.Theory, songio.FromLines(m.Context.Lines), by)
	meta, err := songio.ReadMeta(m.Context.Theory, transposed, true)
	if err != nil {
		return message.UpdateStatusError(err)
	}

	return message.UpdateSong(meta, m.Context.Lines)
}
