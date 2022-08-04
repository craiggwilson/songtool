package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/craiggwilson/songtool/pkg/songio"
)

type ViewCmd struct {
	songCmd
}

func (cmd *ViewCmd) Run(cfg *Config) error {
	defer cmd.ensurePath().Close()

	song := cmd.openSong(cfg)

	mem := songio.NewMemory(song)
	meta, err := songio.ReadMeta(cfg.Theory, mem, false)
	if err != nil {
		return err
	}
	if len(meta.Title) == 0 {
		meta.Title = cmd.Path.Name()
	}

	mem.Rewind()

	p := tea.NewProgram(
		songViewModel{
			cfg:     cfg,
			meta:    meta,
			content: mem,
		},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	return p.Start()
}

type songViewModel struct {
	cfg     *Config
	meta    songio.Meta
	content *songio.Memory

	ready    bool
	viewport viewport.Model
}

func (m songViewModel) Init() tea.Cmd {
	return nil
}

func (m songViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.ready = true
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		m.viewport.SetContent(m.contentView(msg.Width))
		m.content.Rewind()
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m songViewModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m songViewModel) headerView() string {
	header := m.meta.Title
	if m.meta.Key != nil {
		header += fmt.Sprintf(" [%s]", chordStyle.Render(m.meta.Key.Name))
	}

	title := headerStyle.Render(titleStyle.Render(header))
	line := headerFooterBoundaryStyle.Render(strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title))))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m songViewModel) footerView() string {
	info := footerStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := headerFooterBoundaryStyle.Render(strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info))))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m songViewModel) contentView(viewPortWidth int) string {
	sections := m.buildSections()
	maxSectionWidth := 0

	renderedSections := make([]string, len(sections))
	for i, section := range sections {
		rs := section.render()
		maxSectionWidth = max(maxSectionWidth, lipgloss.Width(rs)+colStyle.GetHorizontalFrameSize())
		renderedSections[i] = rs
	}

	numCols := min(viewPortWidth/maxSectionWidth, m.cfg.Styles.MaxColumns)
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

	content := lipgloss.JoinHorizontal(lipgloss.Top, renderedColumns...)

	return lipgloss.PlaceHorizontal(viewPortWidth, lipgloss.Left, content)
}

func (m songViewModel) buildSections() []section {
	var sections []section
	var currentSection section
	for line, ok := m.content.Next(); ok; line, ok = m.content.Next() {
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
