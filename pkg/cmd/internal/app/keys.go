package app

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/status"
)

var defaultKeyMap = keyMap{
	Normal: keyMapNormal{
		CommandMode:    key.NewBinding(key.WithKeys(":"), key.WithHelp(":", "command mode")),
		Help:           key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
		Quit:           key.NewBinding(key.WithKeys("q", "ctlr+c", "esc"), key.WithHelp("q", "quit")),
		Transpose:      key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "transpose")),
		TransposeDown1: key.NewBinding(key.WithKeys("h", "left"), key.WithHelp("←/h", "stepdown")),
		TransposeUp1:   key.NewBinding(key.WithKeys("l", "right"), key.WithHelp("→/l", "stepup")),
	},
	Command: status.DefaultKeyMap(),
	Song:    viewport.DefaultKeyMap(),
}

type keyMap struct {
	Mode Mode

	Normal  keyMapNormal
	Command status.KeyMap
	Song    viewport.KeyMap
}

type keyMapNormal struct {
	CommandMode    key.Binding
	Help           key.Binding
	Quit           key.Binding
	Transpose      key.Binding
	TransposeDown1 key.Binding
	TransposeUp1   key.Binding
}

func (km keyMap) ShortHelp(commandMode bool) []key.Binding {
	if commandMode {
		return []key.Binding{km.Command.Accept, km.Command.Clear}
	}

	switch km.Mode {
	case ModeNormal:
		return []key.Binding{km.Normal.Help, km.Normal.Quit, km.Normal.CommandMode, km.Normal.TransposeDown1, km.Normal.TransposeUp1}
	}

	return nil
}

func (km keyMap) FullHelp(commandMode bool) [][]key.Binding {
	if commandMode {
		return [][]key.Binding{{km.Command.Accept, km.Command.Clear}}
	}

	switch km.Mode {
	case ModeNormal:
		return [][]key.Binding{
			{km.Normal.Help, km.Normal.Quit, km.Normal.CommandMode, km.Normal.Transpose, km.Normal.TransposeDown1, km.Normal.TransposeUp1},
			{km.Song.Up, km.Song.Down, km.Song.PageUp, km.Song.PageDown, km.Song.HalfPageUp, km.Song.PageDown},
		}
	}

	return nil
}
