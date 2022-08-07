package app

import "github.com/charmbracelet/bubbles/key"

var defaultKeyMap = keyMap{
	CommandMode:    key.NewBinding(key.WithKeys(":"), key.WithHelp(":", "command mode")),
	Help:           key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	Quit:           key.NewBinding(key.WithKeys("q", "ctlr+c", "esc"), key.WithHelp("q", "quit")),
	Transpose:      key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "transpose")),
	TransposeDown1: key.NewBinding(key.WithKeys("left"), key.WithHelp("←", "stepdown")),
	TransposeUp1:   key.NewBinding(key.WithKeys("right"), key.WithHelp("→", "stepup")),
}

type keyMap struct {
	CommandMode    key.Binding
	Help           key.Binding
	Quit           key.Binding
	Transpose      key.Binding
	TransposeDown1 key.Binding
	TransposeUp1   key.Binding
}

func (km keyMap) ShortHelp() []key.Binding {
	return []key.Binding{km.Help, km.Quit, km.CommandMode, km.TransposeDown1, km.TransposeUp1}
}

func (km keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{km.Help, km.Quit, km.CommandMode},
		{km.TransposeDown1, km.TransposeUp1},
	}
}
