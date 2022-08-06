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

func NewApp(cfg *config.Config, cmds ...tea.Cmd) appModel {

	return appModel{
		cfg:        cfg,
		initCmds:   cmds,
		header:     newHeaderViewModel(),
		song:       newSongViewModel(cfg.Styles.MaxColumns),
		commandBar: textinput.New(),
	}
}

type appModel struct {
	cfg      *config.Config
	initCmds []tea.Cmd

	ready       bool
	commandMode bool
	err         error

	header     headerViewModel
	song       songViewModel
	commandBar textinput.Model
}

func (m appModel) Init() tea.Cmd {
	initStyles(&m.cfg.Styles)

	cmds := m.initCmds

	return tea.Batch(cmds...)
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case StatusErrorMsg:
		m.err = msg
	case TransposeMsg:
		transposed := songio.Transpose(m.cfg.Theory, songio.FromLines(m.song.Lines), msg.Interval)
		meta, err := songio.ReadMeta(m.cfg.Theory, transposed, true)
		if err != nil {
			cmds = append(cmds, StatusError(err))
		} else {
			cmds = append(cmds, UpdateSong(meta, m.song.Lines))
		}
	case UpdateSongMsg:
		m.header.Meta = &msg.Meta
		m.song.Lines = msg.Lines
	case tea.KeyMsg:
		if m.commandMode {
			switch {
			case key.Matches(msg, defaultKeyMap.command.Accept):
				cmds = append(cmds, runCommand(m.commandContext(), m.commandBar.Value()))
				m.commandBar.Reset()
				m.commandMode = false
			case key.Matches(msg, defaultKeyMap.command.Clear):
				m.commandMode = false
				m.commandBar.Reset()
				cmds = append(cmds, StatusError(nil))
			default:
				m.commandBar, cmd = m.commandBar.Update(msg)
				cmds = append(cmds, cmd)
			}
		} else {
			switch {
			case key.Matches(msg, defaultKeyMap.normal.CommandMode):
				m.commandMode = true
				cmds = append(cmds, StatusError(nil))
				cmds = append(cmds, m.commandBar.Focus())
			case key.Matches(msg, defaultKeyMap.normal.Quit):
				if m.err != nil {
					m.err = nil
					break
				}

				return m, tea.Quit
			case key.Matches(msg, defaultKeyMap.normal.Transpose):
				m.commandMode = true
				m.err = nil
				m.commandBar.SetValue("transpose ")
				cmds = append(cmds, m.commandBar.Focus())
			case key.Matches(msg, defaultKeyMap.normal.TransposeDown1):
				cmds = append(cmds, runCommand(m.commandContext(), "transpose -- -1"))
			case key.Matches(msg, defaultKeyMap.normal.TransposeUp1):
				cmds = append(cmds, runCommand(m.commandContext(), "transpose 1"))
			}
		}

	case tea.WindowSizeMsg:
		m.ready = true

		headerHeight := lipgloss.Height(m.header.View())
		footerHeight := lipgloss.Height(m.footerView())
		commandBarHeight := lipgloss.Height(m.commandBar.View())
		verticalMarginHeight := headerHeight + footerHeight + commandBarHeight

		m.header.Width = msg.Width
		m.song.Width = msg.Width
		m.song.Height = msg.Height - verticalMarginHeight
		m.commandBar.Width = msg.Width
	}

	m.header, cmd = m.header.Update(msg)
	cmds = append(cmds, cmd)

	m.song, cmd = m.song.Update(msg)
	cmds = append(cmds, cmd)

	// todo footer

	return m, tea.Batch(cmds...)
}

func (m appModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	commandBarStr := "\n"
	if m.commandMode {
		commandBarStr += m.commandBar.View()
	} else if m.err != nil {
		commandBarStr += lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render(m.err.Error())
	}

	return fmt.Sprintf("%s\n%s\n%s%s", m.header.View(), m.song.View(), m.footerView(), commandBarStr)
}

func (m *appModel) commandContext() *commandContext {
	return &commandContext{
		Theory: m.cfg.Theory,
		Meta:   m.header.Meta,
		Lines:  m.song.Lines,
	}
}

func (m *appModel) footerView() string {
	info := footerStyle.Render(fmt.Sprintf("%3.f%%", m.song.ScrollPercent()*100))
	line := headerFooterBoundaryStyle.Render(strings.Repeat("─", max(0, m.song.Width-lipgloss.Width(info))))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
