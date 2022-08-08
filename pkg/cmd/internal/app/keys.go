package app

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/explorer"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/song"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/status"
)

var defaultKeyMap = func() *keyMap {
	km := keyMap{
		Global: &globalKeyMap{
			CommandMode: key.NewBinding(key.WithKeys(":"), key.WithHelp(":", "command mode")),
			Explorer:    key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "explorer mode")),
			Song:        key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "song mode")),
			Help:        key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
			Quit:        key.NewBinding(key.WithKeys("q", "esc", "ctlr+c"), key.WithHelp("q/esc", "quit")),
		},
		Explorer: explorer.DefaultKeyMap(),
		Song:     song.DefaultKeyMap(),
		Command:  status.DefaultKeyMap(),
	}

	km.Command.SetEnabled(false)
	return &km
}()

type keyMap struct {
	Global   *globalKeyMap
	Explorer *explorer.KeyMap
	Song     *song.KeyMap
	Command  *status.KeyMap
}

type globalKeyMap struct {
	CommandMode key.Binding
	Explorer    key.Binding
	Song        key.Binding
	Help        key.Binding
	Quit        key.Binding
}

func (km *globalKeyMap) SetEnabled(enabled bool) {
	km.CommandMode.SetEnabled(enabled)
	km.Explorer.SetEnabled(enabled)
	km.Song.SetEnabled(enabled)
	km.Help.SetEnabled(enabled)
	km.Quit.SetEnabled(enabled)
}

func (km *keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		km.Command.Accept,
		km.Command.Clear,
		km.Global.Help,
		km.Global.Quit,
		km.Global.CommandMode,
		km.Global.Explorer,
		km.Global.Song,
		km.Song.TransposeDown1,
		km.Song.TransposeUp1,
	}
}

func (km *keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{km.Command.Accept, km.Command.Clear},
		{km.Global.Help, km.Global.Quit, km.Global.CommandMode, km.Global.Explorer, km.Global.Song},
		{km.Song.Transpose, km.Song.TransposeDown1, km.Song.TransposeUp1},
		{km.Song.Up, km.Song.Down, km.Song.PageUp, km.Song.PageDown, km.Song.HalfPageUp, km.Song.PageDown},
	}
}
