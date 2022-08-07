package status

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Accept key.Binding
	Clear  key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Accept: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "apply")),
		Clear:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "clear and exit")),
	}
}

type HelpKeyMap interface {
	ShortHelp(commandMode bool) []key.Binding
	FullHelp(commandMode bool) [][]key.Binding
}
