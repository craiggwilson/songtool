package songtext

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/message"
	"github.com/craiggwilson/songtool/pkg/songio"
)

var (
	colStyle = lipgloss.NewStyle().MarginLeft(2)
)

func New() Model {
	return Model{
		viewport: viewport.New(0, 0),
	}
}

type Model struct {
	ChordStyle       lipgloss.Style
	LyricsStyle      lipgloss.Style
	SectionNameStyle lipgloss.Style
	KeyMap           KeyMap
	Height           int
	MaxColumns       int
	Width            int

	Lines []songio.Line

	viewport viewport.Model
}

func (m Model) ScrollPercent() float64 {
	return m.viewport.ScrollPercent()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch tmsg := msg.(type) {
	case message.UpdateSongMsg:
		m.Lines = tmsg.Lines
		m.viewport.SetContent(m.contentView())
		return m, nil
	case message.InvalidateMsg:
		m.viewport.Height = m.Height
		m.viewport.Width = m.Width
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.viewport.View()
}

func (m Model) contentView() string {
	if len(m.Lines) == 0 {
		return ""
	}

	sections := m.buildSections()
	maxSectionWidth := 0

	renderedSections := make([]string, len(sections))
	for i, section := range sections {
		rs := section.render(m.SectionNameStyle)
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

func (m Model) buildSections() []section {
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
			currentSection.lines = append(currentSection.lines, m.LyricsStyle.Render(tl.Text))
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
				row += m.ChordStyle.Render(chordName)
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

func (s section) render(sectionNameStyle lipgloss.Style) string {
	var sectionBuilder strings.Builder
	sectionBuilder.WriteString(sectionNameStyle.Render(s.name) + "\n")
	for j, line := range s.lines {
		if j != 0 {
			sectionBuilder.WriteByte('\n')
		}
		sectionBuilder.WriteString(line)
	}

	return sectionBuilder.String()
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
