package models

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	normal  keyMapNormalMode
	command keyMapCommandMode
}

type keyMapNormalMode struct {
	CommandMode    key.Binding
	Quit           key.Binding
	TransposeDown1 key.Binding
	TransposeUp1   key.Binding
}

type keyMapCommandMode struct {
	Accept key.Binding
	Clear  key.Binding
}

var defaultKeyMap = keyMap{
	normal: keyMapNormalMode{
		CommandMode:    key.NewBinding(key.WithKeys(":")),
		Quit:           key.NewBinding(key.WithKeys("esc", "q")),
		TransposeDown1: key.NewBinding(key.WithKeys("ctrl+left")),
		TransposeUp1:   key.NewBinding(key.WithKeys("ctrl+right")),
	},
	command: keyMapCommandMode{
		Accept: key.NewBinding(key.WithKeys("enter")),
		Clear:  key.NewBinding(key.WithKeys("esc")),
	},
}
