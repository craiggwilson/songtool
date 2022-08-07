package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/eval"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/explorer"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/message"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/song"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/status"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/config"
)

func New(cfg *config.Config, cmds ...tea.Cmd) appModel {
	eval := eval.New(cfg.Theory)

	song := song.New(cfg)
	song.KeyMap = defaultKeyMap.Song

	status := status.New()
	status.KeyMap = defaultKeyMap.Command
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
	mode  mode

	eval     eval.Model
	song     song.Model
	explorer explorer.Model
	status   status.Model

	helpModeFull bool
	hasStatus    bool

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

	switch tmsg := msg.(type) {
	case message.EnterCommandModeMsg:
		m.mode |= modeCommand
		defaultKeyMap.Command.SetEnabled(true)
		defaultKeyMap.Explorer.SetEnabled(false)
		defaultKeyMap.Global.SetEnabled(false)
		defaultKeyMap.Song.SetEnabled(false)
	case message.ExitCommandModeMsg:
		m.mode ^= modeCommand
		defaultKeyMap.Command.SetEnabled(false)
		defaultKeyMap.Explorer.SetEnabled(m.mode.IsSongMode())
		defaultKeyMap.Global.SetEnabled(true)
		defaultKeyMap.Song.SetEnabled(m.mode.IsSongMode())
	case message.EnterExplorerModeMsg:
		cmdMode := m.mode.IsCommandMode()
		m.mode = modeExplorer
		if cmdMode {
			m.mode |= modeCommand
		}
		defaultKeyMap.Global.Explorer.SetEnabled(false)
		defaultKeyMap.Global.Song.SetEnabled(true)
		defaultKeyMap.Explorer.SetEnabled(true)
		defaultKeyMap.Song.SetEnabled(false)
		cmds = append(cmds, message.Invalidate())
	case message.EnterSongModeMsg:
		cmdMode := m.mode.IsCommandMode()
		m.mode = modeSong
		if cmdMode {
			m.mode |= modeCommand
		}
		defaultKeyMap.Global.Explorer.SetEnabled(true)
		defaultKeyMap.Global.Song.SetEnabled(false)
		defaultKeyMap.Explorer.SetEnabled(false)
		defaultKeyMap.Song.SetEnabled(true)
		cmds = append(cmds, message.Invalidate())
	case message.InvalidateMsg:
		m.song.Height = m.height - lipgloss.Height(m.status.View())
	case message.UpdateStatusMsg:
		m.hasStatus = tmsg.Info != "" || tmsg.Err != nil
	case tea.KeyMsg:
		if m.mode.IsCommandMode() {
			m.status, cmd = m.status.Update(msg)
			return m, cmd
		} else {
			switch {
			case key.Matches(tmsg, defaultKeyMap.Global.CommandMode):
				return m, message.EnterCommandMode("")
			case key.Matches(tmsg, defaultKeyMap.Global.Explorer):
				return m, message.EnterExplorerMode()
			case key.Matches(tmsg, defaultKeyMap.Global.Song):
				return m, message.EnterSongMode()
			case key.Matches(tmsg, defaultKeyMap.Global.Help):
				m.helpModeFull = !m.helpModeFull
				return m, message.ChangeHelpMode(m.helpModeFull)
			case key.Matches(tmsg, defaultKeyMap.Global.Quit):
				if m.hasStatus {
					return m, message.ClearStatus()
				}

				return m, tea.Quit
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
	cmds = append(cmds, cmd)

	if m.mode.IsSongMode() {
		m.song, cmd = m.song.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.mode.IsExplorerMode() {
		m.explorer, cmd = m.explorer.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.status, cmd = m.status.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m appModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	if m.mode.IsSongMode() {
		return fmt.Sprintf("%s\n%s", m.song.View(), m.status.View())
	} else if m.mode.IsExplorerMode() {
		return fmt.Sprintf("%s\n%s", m.explorer.View(), m.status.View())
	}

	return "\n  Problems..."
}
