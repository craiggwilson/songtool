package status

import "github.com/charmbracelet/bubbles/key"

func DefaultKeyMap() *KeyMap {
	return &KeyMap{
		Accept: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "apply")),
		Clear:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "clear and exit")),
	}
}

type KeyMap struct {
	Accept key.Binding
	Clear  key.Binding
}

func (km *KeyMap) SetEnabled(enabled bool) {
	km.Accept.SetEnabled(enabled)
	km.Clear.SetEnabled(enabled)
}
