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
	CommandMode  bool
	FullHelpMode bool
	HelpKeyMap   help.KeyMap
	KeyMap       KeyMap
	Width        int

	Info string
	Err  error

	command textinput.Model
	help    help.Model
}

func (m *Model) Focus() tea.Cmd {
	return m.command.Focus()
}

func (m *Model) SetValue(value string) {
	m.command.SetValue(value)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	appendCmd := func(cmd tea.Cmd) {
		cmds = append(cmds, cmd)
	}

	switch tmsg := msg.(type) {
	case message.UpdateStatusMsg:
		m.Err = tmsg.Err
	case tea.KeyMsg:
		if m.CommandMode {
			switch {
			case key.Matches(tmsg, m.KeyMap.Accept):
				appendCmd(message.Eval(m.command.Value()))
				m.command.Reset()
				m.CommandMode = false
				msg = nil
			case key.Matches(tmsg, m.KeyMap.Clear):
				appendCmd(message.UpdateStatusError(nil))
				m.command.Reset()
				m.CommandMode = false
				msg = nil
			default:
				m.command, cmd = m.command.Update(msg)
				appendCmd(cmd)
				msg = nil
			}
		}
	}

	m.command.Width = m.Width
	m.help.Width = m.Width

	m.help.ShowAll = m.FullHelpMode

	m.command, cmd = m.command.Update(msg)
	appendCmd(cmd)

	m.help, cmd = m.help.Update(msg)
	appendCmd(cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.CommandMode {
		return m.command.View()
	} else if m.Err != nil {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render(m.Err.Error())
	}

	return m.help.View(m.HelpKeyMap)
}
