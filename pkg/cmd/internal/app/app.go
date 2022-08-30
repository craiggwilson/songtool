package app

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/message"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/config"
	"github.com/craiggwilson/teakwood"
	"github.com/craiggwilson/teakwood/adapter"
	"github.com/craiggwilson/teakwood/items"
	"github.com/craiggwilson/teakwood/items/flow"
	"github.com/craiggwilson/teakwood/items/tabs"
	"github.com/craiggwilson/teakwood/named"
	"github.com/craiggwilson/teakwood/stack"
)

func New(cfg *config.Config, cmds ...tea.Cmd) appModel {
	app := appModel{
		cfg:                 cfg,
		mode:                modeSong,
		initCmds:            cmds,
		loadedSongs:         items.NewSlice[*loadedSong](),
		currentSongSections: items.NewSlice[*loadedSongSection](),
	}

	app.songView = stack.New(
		stack.WithOrientation(stack.Vertical),
		stack.WithItems(
			stack.NewAutoItem(
				tabs.New[*loadedSong](
					app.loadedSongs,
					items.RenderFunc[*loadedSong](func(sm *loadedSong) string {
						return sm.meta.Title
					}),
				),
			),
			stack.NewProportionalItem(1,
				named.New("content", flow.New[*loadedSongSection](
					app.currentSongSections,
					items.RenderFunc[*loadedSongSection](func(sm *loadedSongSection) string {
						return "here"
					}),
					flow.WithHorizontalAlignment[*loadedSongSection](lipgloss.Center),
					flow.WithOrientation[*loadedSongSection](flow.Vertical),
					flow.WithVerticalAlignment[*loadedSongSection](lipgloss.Center),
				)),
			),
			stack.NewAutoItem(
				named.New("help", adapter.New(
					help.New(),
					adapter.WithUpdateBounds(func(m help.Model, bounds teakwood.Rectangle) help.Model {
						m.Width = bounds.Width
						return m
					}),
					adapter.WithView(func(m help.Model) string {
						return m.View(defaultKeyMap)
					}),
				)),
			),
		),
	)

	return app
}

type appModel struct {
	cfg      *config.Config
	initCmds []tea.Cmd

	ready bool
	mode  mode

	loadedSongs         *items.Slice[*loadedSong]
	currentSongIdx      int
	currentSongSections items.Source[*loadedSongSection]

	songView stack.Model

	hasStatus bool
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
	case SongOpenedMsg:
		m.loadedSongs.Add(&loadedSong{
			path:  tmsg.Path,
			meta:  &tmsg.Meta,
			lines: tmsg.Lines,
		})
		return m, nil
	case OpenSongMsg:
		return m, m.openSong(tmsg.Path)
	case tea.KeyMsg:
		switch {
		case m.mode.IsCommandMode():
			return m, cmd
		case key.Matches(tmsg, defaultKeyMap.Global.CommandMode):
			return m, message.EnterCommandMode("")
		case key.Matches(tmsg, defaultKeyMap.Global.Explorer):
			return m, message.EnterExplorerMode()
		case key.Matches(tmsg, defaultKeyMap.Global.Song):
			return m, message.EnterSongMode()
		case key.Matches(tmsg, defaultKeyMap.Global.Help):
			return m, named.Update("help", func(a adapter.Model[help.Model], m2 tea.Msg) (tea.Model, tea.Cmd) {
				a.UpdateAdaptee(func(h help.Model) help.Model {
					h.ShowAll = !h.ShowAll
					return h
				})
				return a, nil
			})
		case key.Matches(tmsg, defaultKeyMap.Global.Quit):
			if m.hasStatus {
				return m, message.ClearStatus()
			}

			return m, tea.Quit
		case m.mode.IsExplorerMode():
			return m, cmd
		case m.mode.IsSongMode():
			return m, cmd
		}
	case tea.WindowSizeMsg:
		m.ready = true
		v := m.songView.UpdateBounds(teakwood.NewRectangle(0, 0, tmsg.Width-1, tmsg.Height-1))
		m.songView = v.(stack.Model)
	}

	newSongView, cmd := m.songView.Update(msg)
	m.songView = newSongView.(stack.Model)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m appModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	if m.mode.IsSongMode() {
		return m.songView.View()
	}

	return "\n  Problems..."
}

func (m *appModel) updateKeyBindings() {
	defaultKeyMap.Command.SetEnabled(m.mode.IsCommandMode())

	if m.mode.IsCommandMode() {
		defaultKeyMap.Global.SetEnabled(false)

		defaultKeyMap.Explorer.SetEnabled(false)
		defaultKeyMap.Song.SetEnabled(false)
	} else {
		defaultKeyMap.Global.SetEnabled(true)

		defaultKeyMap.Global.Explorer.SetEnabled(!m.mode.IsExplorerMode())
		defaultKeyMap.Global.Song.SetEnabled(!m.mode.IsSongMode())
		defaultKeyMap.Explorer.SetEnabled(m.mode.IsExplorerMode())
		defaultKeyMap.Song.SetEnabled(m.mode.IsSongMode())
	}
}
