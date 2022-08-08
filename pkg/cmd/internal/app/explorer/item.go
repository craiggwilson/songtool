package explorer

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	Path  string
	Title string
}

func (i item) FilterValue() string {
	return i.Title
}

func newItemDelegate() itemDelegate {
	return itemDelegate{}
}

type itemDelegate struct{}

func (d itemDelegate) Height() int {
	return 1
}

func (d itemDelegate) Spacing() int {
	return 0
}

func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Path)

	if index == m.Index() {
		str = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render(str)
	}

	fmt.Fprintf(w, str)
}
