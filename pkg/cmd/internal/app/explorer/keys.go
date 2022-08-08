package explorer

import "github.com/charmbracelet/bubbles/key"

func DefaultKeyMap() *KeyMap {
	return &KeyMap{
		Select:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select song")),
		MoveRight: key.NewBinding(key.WithKeys("right", "k"), key.WithHelp("right", "move right")),
		MoveLeft:  key.NewBinding(key.WithKeys("left", "j"), key.WithHelp("left", "move left")),
		MoveUp:    key.NewBinding(key.WithKeys("up", "h"), key.WithHelp("up", "move up")),
		MoveDown:  key.NewBinding(key.WithKeys("down", "l"), key.WithHelp("down", "move down")),
	}
}

type KeyMap struct {
	Select    key.Binding
	MoveRight key.Binding
	MoveLeft  key.Binding
	MoveUp    key.Binding
	MoveDown  key.Binding
}

func (km *KeyMap) SetEnabled(enabled bool) {
	km.Select.SetEnabled(enabled)
}
