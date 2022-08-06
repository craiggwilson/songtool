package models

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/config"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/models/bubbles/header"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/models/message"
	"github.com/craiggwilson/songtool/pkg/songio"
	"github.com/craiggwilson/songtool/pkg/theory/interval"
)

func NewApp(cfg *config.Config, cmds ...tea.Cmd) appModel {
	initStyles(&cfg.Styles)

	header := header.New()
	header.BorderColor = cfg.Styles.BoundaryColor.Color()
	header.KeyStyle = cfg.Styles.Chord.Style()
	header.TitleStyle = cfg.Styles.Title.Style()

	return appModel{
		cfg:        cfg,
		initCmds:   cmds,
		header:     header,
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
	meta        *songio.Meta
	lines       []songio.Line

	header     header.Header
	song       songViewModel
	commandBar textinput.Model
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

	switch msg := msg.(type) {
	case message.LoadSongMsg:
		appendCmd(m.loadSong(msg.Path))
	case message.TransposeSongMsg:
		appendCmd(m.transposeSong(msg.Interval))
	case message.UpdateSongMsg:
		m.meta = &msg.Meta
		m.lines = msg.Lines
		m.header.Meta = m.meta
		m.song.Lines = m.lines
	case message.UpdateStatusMsg:
		m.err = msg.Err
	case tea.KeyMsg:
		if m.commandMode {
			switch {
			case key.Matches(msg, defaultKeyMap.command.Accept):
				appendCmd(runCommand(m.commandContext(), m.commandBar.Value()))
				m.commandBar.Reset()
				m.commandMode = false
			case key.Matches(msg, defaultKeyMap.command.Clear):
				m.commandMode = false
				m.commandBar.Reset()
				appendCmd(message.UpdateStatusError(nil))
			default:
				m.commandBar, cmd = m.commandBar.Update(msg)
				appendCmd(cmd)
			}
		} else {
			switch {
			case key.Matches(msg, defaultKeyMap.normal.CommandMode):
				m.commandMode = true
				appendCmd(message.UpdateStatusError(nil))
				appendCmd(m.commandBar.Focus())
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
				appendCmd(m.commandBar.Focus())
			case key.Matches(msg, defaultKeyMap.normal.TransposeDown1):
				appendCmd(runCommand(m.commandContext(), "transpose -- -1"))
			case key.Matches(msg, defaultKeyMap.normal.TransposeUp1):
				appendCmd(runCommand(m.commandContext(), "transpose 1"))
			}
		}

	case tea.WindowSizeMsg:
		m.ready = true

		m.header.Width = msg.Width
		m.song.Width = msg.Width
		m.commandBar.Width = msg.Width

		headerHeight := lipgloss.Height(m.header.View())
		commandBarHeight := lipgloss.Height(m.commandBarView())
		footerHeight := lipgloss.Height(m.footerView())

		verticalMarginHeight := headerHeight + footerHeight + commandBarHeight + 1

		m.song.Height = msg.Height - verticalMarginHeight
	}

	m.header, cmd = m.header.Update(msg)
	appendCmd(cmd)

	m.song, cmd = m.song.Update(msg)
	appendCmd(cmd)

	// todo footer

	return m, tea.Batch(cmds...)
}

func (m appModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	return fmt.Sprintf("%s\n%s\n%s%s", m.header.View(), m.song.View(), m.footerView(), m.commandBarView())
}

func (m *appModel) commandContext() *commandContext {
	return &commandContext{
		Theory: m.cfg.Theory,
		Meta:   m.header.Meta,
		Lines:  m.song.Lines,
	}
}

func (m appModel) commandBarView() string {
	commandBarStr := "\n"
	if m.commandMode {
		commandBarStr += m.commandBar.View()
	} else if m.err != nil {
		commandBarStr += lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render(m.err.Error())
	}

	return commandBarStr
}

func (m appModel) footerView() string {
	info := footerStyle.Render(fmt.Sprintf("%3.f%%", m.song.ScrollPercent()*100))
	line := headerFooterBoundaryStyle.Render(strings.Repeat("─", max(0, m.song.Width-lipgloss.Width(info))))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m appModel) loadSong(path string) tea.Cmd {
	var f *os.File
	var err error
	switch path {
	case "":
		return message.UpdateStatusError(fmt.Errorf("no file to load"))
	case "-":
		f = os.Stdin
	default:
		f, err = os.Open(path)
		if err != nil {
			return message.UpdateStatusError(err)
		}
	}
	defer f.Close()

	rdr := songio.ReadChordsOverLyrics(m.cfg.Theory, m.cfg.Theory, f)
	lines, err := songio.ReadAllLines(rdr)
	if err != nil {
		return message.UpdateStatusError(err)
	}

	meta, err := songio.ReadMeta(m.cfg.Theory, songio.FromLines(lines), true)
	if err != nil {
		return message.UpdateStatusError(err)
	}

	if meta.Title == "" {
		meta.Title = path
	}

	return message.UpdateSong(meta, lines)
}

func (m appModel) transposeSong(by interval.Interval) tea.Cmd {
	transposed := songio.Transpose(m.cfg.Theory, songio.FromLines(m.lines), by)
	meta, err := songio.ReadMeta(m.cfg.Theory, transposed, true)
	if err != nil {
		return message.UpdateStatusError(err)
	}

	return message.UpdateSong(meta, m.song.Lines)
}
