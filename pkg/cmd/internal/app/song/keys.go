package song

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
)

func DefaultKeyMap() KeyMap {
	return KeyMap{
		KeyMap:         viewport.DefaultKeyMap(),
		Transpose:      key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "transpose")),
		TransposeDown1: key.NewBinding(key.WithKeys("h", "left"), key.WithHelp("←/h", "transpose down")),
		TransposeUp1:   key.NewBinding(key.WithKeys("l", "right"), key.WithHelp("→/l", "transpose up")),
	}
}

type KeyMap struct {
	viewport.KeyMap

	Transpose      key.Binding
	TransposeDown1 key.Binding
	TransposeUp1   key.Binding
}

func (km *KeyMap) SetEnabled(enabled bool) {
	km.Transpose.SetEnabled(enabled)
	km.TransposeDown1.SetEnabled(enabled)
	km.TransposeUp1.SetEnabled(enabled)

	km.KeyMap.Down.SetEnabled(enabled)
	km.KeyMap.Up.SetEnabled(enabled)
	km.KeyMap.HalfPageDown.SetEnabled(enabled)
	km.KeyMap.HalfPageUp.SetEnabled(enabled)
	km.KeyMap.PageDown.SetEnabled(enabled)
	km.KeyMap.PageUp.SetEnabled(enabled)
}
