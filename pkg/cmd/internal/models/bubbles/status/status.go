package status

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/models/message"
)

func New() Model {
	return Model{
		KeyMap:  DefaultKeyMap(),
		command: textinput.New(),
	}
}

type Model struct {
	CommandMode bool
	KeyMap      KeyMap
	Width       int

	Info string
	Err  error

	command textinput.Model
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
			case key.Matches(tmsg, m.KeyMap.Clear):
				appendCmd(message.UpdateStatusError(nil))
				m.command.Reset()
				m.CommandMode = false
			default:
				m.command, cmd = m.command.Update(msg)
				appendCmd(cmd)
				msg = nil
			}
		}
	}

	m.command, cmd = m.command.Update(msg)
	appendCmd(cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	v := "\n"
	if m.CommandMode {
		v += m.command.View()
	} else if m.Err != nil {
		v += lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render(m.Err.Error())
	}

	return v
}
