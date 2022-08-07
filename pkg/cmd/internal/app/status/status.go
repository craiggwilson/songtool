package status

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/message"
)

func New() Model {
	return Model{
		KeyMap:  DefaultKeyMap(),
		command: textinput.New(),
		help:    help.New(),
	}
}

type Model struct {
	HelpKeyMap HelpKeyMap
	KeyMap     KeyMap
	Width      int

	info string
	err  error

	commandMode bool
	command     textinput.Model
	help        help.Model
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.command.Width = m.Width
	m.help.Width = m.Width

	switch tmsg := msg.(type) {
	case message.ChangeHelpModeMsg:
		m.help.ShowAll = tmsg.Full
		m.help, cmd = m.help.Update(msg)
		return m, tea.Batch(cmd, message.Invalidate())
	case message.EnterCommandModeMsg:
		m.err = nil
		m.commandMode = true
		if tmsg.Value != "" {
			m.command.SetValue(tmsg.Value)
		}
		cmd = m.command.Focus()
		return m, tea.Batch(message.Invalidate(), cmd)
	case message.ExitCommandModeMsg:
		m.commandMode = false
		return m, message.Invalidate()
	case message.UpdateStatusMsg:
		m.info = tmsg.Info
		m.err = tmsg.Err
		return m, message.Invalidate()
	case tea.KeyMsg:
		if m.commandMode {
			switch {
			case key.Matches(tmsg, m.KeyMap.Accept):
				value := m.command.Value()
				m.command.Reset()
				return m, tea.Batch(
					message.ExitCommandMode(),
					message.Eval(value),
				)
			case key.Matches(tmsg, m.KeyMap.Clear):
				m.command.Reset()
				m.err = nil
				return m, message.ExitCommandMode()
			default:
				m.command, cmd = m.command.Update(msg)
				return m, cmd
			}
		}
	}

	if _, ok := msg.(tea.KeyMsg); !ok {
		m.command, cmd = m.command.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.help, cmd = m.help.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	v := ""
	if m.commandMode {
		v += m.command.View() + "\n"
	} else if m.err != nil {
		v += lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render(m.err.Error()) + "\n"
	}

	if m.help.ShowAll {
		return v + m.help.FullHelpView(m.HelpKeyMap.FullHelp(m.commandMode))
	}

	return v + m.help.ShortHelpView(m.HelpKeyMap.ShortHelp(m.commandMode))
}
