package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
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
	errStr      string

	viewport   songViewPortModel
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

	m.viewport.Lines = lines

	return nil
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
		err  error
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.commandMode {
			switch {
			case key.Matches(msg, defaultKeyMap.command.Accept):
				cmd, err := runCommand(&m, m.commandBar.Value())
				cmds = append(cmds, cmd)
				if err != nil {
					m.errStr = err.Error()
				} else {
					m.commandBar.SetValue("")
				}
				m.commandMode = false
			case key.Matches(msg, defaultKeyMap.command.Clear):
				m.commandMode = false
				m.commandBar.SetValue("")
			default:
				m.commandBar, cmd = m.commandBar.Update(msg)
				cmds = append(cmds, cmd)
			}
		} else {
			switch {
			case key.Matches(msg, defaultKeyMap.normal.CommandMode):
				m.commandMode = true
				m.errStr = ""
				m.commandBar.Focus()
			case key.Matches(msg, defaultKeyMap.normal.Quit):
				if m.errStr != "" {
					m.errStr = ""
					break
				}

				return m, tea.Quit
			case key.Matches(msg, defaultKeyMap.normal.Transpose):
				m.commandMode = true
				m.errStr = ""
				m.commandBar.SetValue("transpose ")
				m.commandBar.Focus()
			case key.Matches(msg, defaultKeyMap.normal.TransposeDown1):
				cmd, err = runCommand(&m, "transpose -- -1")
				cmds = append(cmds, cmd)
				if err != nil {
					m.errStr = err.Error()
				}
			case key.Matches(msg, defaultKeyMap.normal.TransposeUp1):
				cmd, err = runCommand(&m, "transpose 1")
				cmds = append(cmds, cmd)
				if err != nil {
					m.errStr = err.Error()
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
			m.viewport = newSongViewPort(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.MaxColumns = m.cfg.Styles.MaxColumns
			m.viewport.Lines = m.lines
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
	} else if len(m.errStr) > 0 {
		commandBarStr += lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render(m.errStr)
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
