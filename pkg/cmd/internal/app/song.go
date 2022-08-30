package app

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/message"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/config"
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
	sections []string
}

func (sm *loadedSong) Sections(styles config.Styles) []string {
	if sm.sections != nil {
		return sm.sections
	}

	var currentName string
	var currentLines []string
	for _, line := range sm.lines {
		switch tl := line.(type) {
		case *songio.SectionStartDirectiveLine:
			currentName = tl.Name
		case *songio.SectionEndDirectiveLine:
			sm.sections = append(sm.sections, sm.renderSection(styles, currentName, currentLines))
		case *songio.TextLine:
			currentLines = append(currentLines, styles.Lyrics.Render(tl.Text))
		case *songio.ChordLine:
			row := ""
			currentOffset := 0
			for _, chordOffset := range tl.Chords {
				offsetDiff := chordOffset.Offset - currentOffset
				if offsetDiff > 0 {
					row += strings.Repeat(" ", offsetDiff)
					currentOffset += offsetDiff
				}

				chordName := chordOffset.Chord.Name
				row += styles.Chord.Render(chordName)
				currentOffset += len(chordName)
			}

			currentLines = append(currentLines, row)
		case *songio.EmptyLine:
			currentLines = append(currentLines, "")
		}
	}

	return sm.sections
}

func (sm *loadedSong) renderSection(styles config.Styles, name string, lines []string) string {
	var sectionBuilder strings.Builder
	sectionBuilder.WriteString(styles.SectionName.Render(name) + "\n")
	for j, line := range lines {
		if j != 0 {
			sectionBuilder.WriteByte('\n')
		}
		sectionBuilder.WriteString(line)
	}

	return sectionBuilder.String()
}
