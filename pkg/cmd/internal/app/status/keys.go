package status

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Accept key.Binding
	Clear  key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Accept: key.NewBinding(key.WithKeys("enter")),
		Clear:  key.NewBinding(key.WithKeys("esc")),
	}
}
