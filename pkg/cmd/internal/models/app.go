package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/config"
	"github.com/craiggwilson/songtool/pkg/songio"
)

func NewApp(cfg *config.Config) appModel {
	return appModel{
		cfg: cfg,
	}
}

type appModel struct {
	cfg   *config.Config
	ready bool

	commandMode bool
	lines       []songio.Line
	meta        songio.Meta

	viewport   SongViewPortModel
	commandBar textinput.Model
}

func (m appModel) Init() tea.Cmd {
	initStyles(&m.cfg.Styles)
	return nil
}

func (m *appModel) SetSong(path string, song songio.Song) error {
	lines, err := songio.ReadAllLines(song)
	if err != nil {
		return err
	}

	meta, err := songio.ReadMeta(m.cfg.Theory, songio.NewMemory(lines), false)
	if err != nil {
		return err
	}
	if len(meta.Title) == 0 {
		meta.Title = path
	}

	m.lines = lines
	m.meta = meta

	m.viewport.SetSongLines(lines)

	return nil
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.commandMode {
			switch msg.Type {
			case tea.KeyEnter:
				switch m.commandBar.Value() {
				case "q":
					return m, tea.Quit
				}
			case tea.KeyEsc, tea.KeyCtrlC:
				m.commandMode = false
				m.commandBar.SetValue("")
			default:
				m.commandBar, cmd = m.commandBar.Update(msg)
				cmds = append(cmds, cmd)
			}
		} else {
			switch msg.Type {
			case tea.KeyEsc, tea.KeyCtrlC:
				return m, tea.Quit
			case tea.KeyRunes:
				switch string(msg.Runes) {
				case "q":
					return m, tea.Quit
				case ":":
					m.commandMode = true
					m.commandBar.Focus()
				}
			}
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		commandBarHeight := lipgloss.Height(m.commandBar.View())
		verticalMarginHeight := headerHeight + footerHeight + commandBarHeight

		m.viewport.Height = msg.Height - verticalMarginHeight

		if !m.ready {
			m.viewport = NewSongViewPort(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.SetSongLines(m.lines)
			m.commandBar = textinput.New()
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
			m.commandBar.Width = msg.Width
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m appModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	commandBarStr := "\n"
	if m.commandMode {
		commandBarStr += m.commandBar.View()
	}

	return fmt.Sprintf("%s\n%s\n%s%s", m.headerView(), m.viewport.View(), m.footerView(), commandBarStr)
}

func (m appModel) headerView() string {
	var title string
	header := m.meta.Title
	if m.meta.Key != nil {
		header += fmt.Sprintf(" [%s]", chordStyle.Render(m.meta.Key.Name))
	}
	title = headerStyle.Render(titleStyle.Render(header))

	line := headerFooterBoundaryStyle.Render(strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title))))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m appModel) footerView() string {
	info := footerStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := headerFooterBoundaryStyle.Render(strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info))))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
