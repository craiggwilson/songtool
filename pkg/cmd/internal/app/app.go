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

const helpName = "help"

func New(cfg *config.Config, cmds ...tea.Cmd) appModel {
	app := appModel{
		cfg:                 cfg,
		keyMap:              defaultKeyMap,
		initCmds:            cmds,
		loadedSongs:         items.NewSlice[*loadedSong](),
		currentSongSections: items.NewSlice[string](),
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
				named.New("content", flow.New[string](
					app.currentSongSections,
					items.RenderFunc[string](func(s string) string {
						return s
					}),
					flow.WithHorizontalAlignment[string](lipgloss.Center),
					flow.WithOrientation[string](flow.Vertical),
					flow.WithVerticalAlignment[string](lipgloss.Center),
					flow.WithStyles[string](flow.Styles{
						Item: lipgloss.NewStyle().Padding(2),
					}),
				)),
			),
			stack.NewAutoItem(
				named.New(helpName, adapter.New(
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
	keyMap   *keyMap

	mode mode

	loadedSongs *items.Slice[*loadedSong]

	currentSongIdx      int
	currentSongSections *items.Slice[string]

	songView tea.Model

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
	case EnterSongModeMsg:
		m.mode = modeSong
		m.updateKeyBindings()
	case OpenSongMsg:
		cmds = append(cmds, m.openSong(tmsg.Path))
	case SongOpenedMsg:
		m.loadedSongs.Add(&loadedSong{
			path:  tmsg.Path,
			meta:  &tmsg.Meta,
			lines: tmsg.Lines,
		})
		m.currentSongIdx = m.loadedSongs.Len() - 1
		m.currentSongSections.Clear()
		m.currentSongSections.Add(m.loadedSongs.Item(m.currentSongIdx).Sections(m.cfg.Styles)...)
	case tea.KeyMsg:
		switch {
		// case m.mode.IsCommandMode():
		// 	return m, cmd
		case key.Matches(tmsg, m.keyMap.Global.CommandMode):
			return m, message.EnterCommandMode("")
		case key.Matches(tmsg, m.keyMap.Global.Explorer):
			return m, message.EnterExplorerMode()
		case key.Matches(tmsg, m.keyMap.Global.Song):
			return m, message.EnterSongMode()
		case key.Matches(tmsg, m.keyMap.Global.Help):
			cmds = append(cmds, named.Update(helpName, func(a adapter.Model[help.Model], _ tea.Msg) (tea.Model, tea.Cmd) {
				a.UpdateAdaptee(func(h help.Model) help.Model {
					h.ShowAll = !h.ShowAll
					return h
				})
				return a, invalidate
			}))
		case key.Matches(tmsg, m.keyMap.Global.Quit):
			if !m.hasStatus {
				return m, tea.Quit
			}

			cmds = append(cmds, message.ClearStatus())
		}
	case tea.WindowSizeMsg:
		m.songView = m.songView.(stack.Model).UpdateBounds(teakwood.NewRectangle(0, 0, tmsg.Width, tmsg.Height))
	}

	m.songView, cmd = m.songView.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m appModel) View() string {
	return m.songView.View()
}

// this is a little hack to handle getting around bounds updates through the stack.
func invalidate() tea.Msg {
	return invalidateMsg{}
}

type invalidateMsg struct{}
