package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/eval"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/message"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/song"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/status"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/config"
)

func New(cfg *config.Config, cmds ...tea.Cmd) appModel {
	eval := eval.New(cfg.Theory)

	song := song.New(cfg)

	status := status.New()
	status.HelpKeyMap = defaultKeyMap

	return appModel{
		cfg:      cfg,
		initCmds: cmds,
		eval:     eval,
		song:     song,
		status:   status,
	}
}

type appModel struct {
	cfg      *config.Config
	initCmds []tea.Cmd

	ready bool

	eval   eval.Model
	song   song.Model
	status status.Model

	height int
	width  int
}

func (m appModel) Init() tea.Cmd {
	return tea.Batch(m.initCmds...)
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	appendCmd := func(cmd tea.Cmd) {
		cmds = append(cmds, cmd)
	}

	switch tmsg := msg.(type) {
	case tea.KeyMsg:
		if m.status.CommandMode {
			m.status, cmd = m.status.Update(msg)
			appendCmd(cmd)
			msg = nil
		} else {
			switch {
			case key.Matches(tmsg, defaultKeyMap.CommandMode):
				m.status.CommandMode = true
				appendCmd(message.UpdateStatusError(nil))
				appendCmd(m.status.Focus())

				msg = nil // handled
			case key.Matches(tmsg, defaultKeyMap.Help):
				m.status.FullHelpMode = !m.status.FullHelpMode
				appendCmd(func() tea.Msg {
					return tea.WindowSizeMsg{
						Width:  m.width,
						Height: m.height,
					}
				})
				msg = nil
			case key.Matches(tmsg, defaultKeyMap.Quit):
				if m.status.Err != nil {
					m.status.Err = nil
					msg = nil
					break
				}

				return m, tea.Quit
			case key.Matches(tmsg, defaultKeyMap.Transpose):
				m.status.Err = nil
				m.status.SetValue("transpose ")
				m.status.CommandMode = true
				appendCmd(m.status.Focus())
				msg = nil
			case key.Matches(tmsg, defaultKeyMap.TransposeDown1):
				appendCmd(message.Eval("transpose -- -1"))
				msg = nil
			case key.Matches(tmsg, defaultKeyMap.TransposeUp1):
				appendCmd(message.Eval("transpose 1"))
				msg = nil
			}
		}

	case tea.WindowSizeMsg:
		m.ready = true
		m.height = tmsg.Height
		m.width = tmsg.Width

		m.song.Width = m.width
		m.status.Width = m.width

		m.song.Height = m.height - lipgloss.Height(m.status.View())
	}

	m.eval, cmd = m.eval.Update(msg)
	appendCmd(cmd)

	m.song, cmd = m.song.Update(msg)
	appendCmd(cmd)

	m.status, cmd = m.status.Update(msg)
	appendCmd(cmd)

	return m, tea.Batch(cmds...)
}

func (m appModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	return fmt.Sprintf("%s\n%s", m.song.View(), m.status.View())
}
