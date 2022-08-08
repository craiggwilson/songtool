package explorer

import "github.com/charmbracelet/bubbles/key"

func DefaultKeyMap() *KeyMap {
	return &KeyMap{
		Select: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select song")),
	}
}

type KeyMap struct {
	Select key.Binding
}

func (km *KeyMap) SetEnabled(enabled bool) {
	km.Select.SetEnabled(enabled)
}
