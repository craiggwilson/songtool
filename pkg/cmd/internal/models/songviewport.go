package models

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/craiggwilson/songtool/pkg/songio"
)

func NewSongViewPort(width, height int) SongViewPortModel {
	return SongViewPortModel{
		Height:     height,
		Width:      width,
		MaxColumns: 3,
		viewport:   viewport.New(width, height),
	}
}

type SongViewPortModel struct {
	Height     int
	MaxColumns int
	Width      int
	Lines      []songio.Line

	ready bool

	viewport viewport.Model

	commandMode bool
	commandBar  textinput.Model
}

func (m SongViewPortModel) ScrollPercent() float64 {
	return m.viewport.ScrollPercent()
}

func (m SongViewPortModel) Update(msg tea.Msg) (SongViewPortModel, tea.Cmd) {
	m.viewport.Width = m.Width
	m.viewport.Height = m.Height

	if m.Lines == nil {
		m.viewport.SetContent("")
	} else {
		m.viewport.SetContent(m.contentView())
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m SongViewPortModel) View() string {
	if m.Lines == nil {
		return "<no song>"
	}

	return m.viewport.View()
}

func (m SongViewPortModel) contentView() string {
	sections := m.buildSections()
	maxSectionWidth := 0

	renderedSections := make([]string, len(sections))
	for i, section := range sections {
		rs := section.render()
		maxSectionWidth = max(maxSectionWidth, lipgloss.Width(rs)+colStyle.GetHorizontalFrameSize())
		renderedSections[i] = rs
	}

	numCols := max(1, min(m.Width/maxSectionWidth, m.MaxColumns))
	colStyle.Width(maxSectionWidth)

	renderedColumns := make([]string, numCols)
	for i := 0; i < len(renderedSections); i += numCols {
		for j := 0; j < numCols && i+j < len(renderedSections); j++ {
			if i >= numCols {
				renderedColumns[j] += "\n\n"
			}
			renderedColumns[j] += renderedSections[i+j]
		}
	}

	for i := 0; i < len(renderedColumns); i++ {
		renderedColumns[i] = colStyle.Render(renderedColumns[i])
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedColumns...)
}

func (m SongViewPortModel) buildSections() []section {
	var sections []section
	var currentSection section
	for _, line := range m.Lines {
		switch tl := line.(type) {
		case *songio.SectionStartDirectiveLine:
			currentSection = section{
				name: tl.Name,
			}
		case *songio.SectionEndDirectiveLine:
			sections = append(sections, currentSection)
		case *songio.TextLine:
			currentSection.lines = append(currentSection.lines, lyricsStyle.Render(tl.Text))
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
				row += chordStyle.Render(chordName)
				currentOffset += len(chordName)
			}

			currentSection.lines = append(currentSection.lines, row)
		case *songio.EmptyLine:
			currentSection.lines = append(currentSection.lines, "")
		}
	}

	return sections
}

type section struct {
	name  string
	lines []string
}

func (s section) render() string {
	var sectionBuilder strings.Builder
	sectionBuilder.WriteString(sectionNameStyle.Render(s.name) + "\n")
	for j, line := range s.lines {
		if j != 0 {
			sectionBuilder.WriteByte('\n')
		}
		sectionBuilder.WriteString(line)
	}

	return sectionStyle.Render(sectionBuilder.String())
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
