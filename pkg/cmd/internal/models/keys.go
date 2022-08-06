package models

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	CommandMode    key.Binding
	Quit           key.Binding
	Transpose      key.Binding
	TransposeDown1 key.Binding
	TransposeUp1   key.Binding
}

var defaultKeyMap = keyMap{
	CommandMode:    key.NewBinding(key.WithKeys(":")),
	Quit:           key.NewBinding(key.WithKeys("esc", "q")),
	Transpose:      key.NewBinding(key.WithKeys("ctrl+t")),
	TransposeDown1: key.NewBinding(key.WithKeys("ctrl+left")),
	TransposeUp1:   key.NewBinding(key.WithKeys("ctrl+right")),
}
