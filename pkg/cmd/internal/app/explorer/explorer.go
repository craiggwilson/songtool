package explorer

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/message"
)

func New() Model {
	items := []list.Item{
		item{Title: "Funny 1"},
		item{Title: "All Creatures"},
	}

	list := list.New(items, newItemDelegate(), 0, 0)
	list.DisableQuitKeybindings()
	list.SetShowHelp(false)
	list.SetShowStatusBar(false)
	list.SetShowTitle(false)

	return Model{
		list: list,
	}
}

type Model struct {
	KeyMap *KeyMap
	Width  int
	Height int

	list list.Model
	path string
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch tmsg := msg.(type) {
	case message.UpdateFilesMsg:
		cmd = m.updateItems(tmsg.Files)
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		switch {
		case key.Matches(tmsg, m.KeyMap.Select):
			item := m.list.SelectedItem().(item)
			return m, tea.Batch(
				message.LoadSong(item.Path),
				message.EnterSongMode(),
			)
		}
	case message.InvalidateMsg:
		m.list.SetWidth(m.Width)
		m.list.SetHeight(m.Height)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.list.View()
}

func (m *Model) updateItems(files []message.FileItem) tea.Cmd {
	items := make([]list.Item, len(files))
	for i := range files {
		items[i] = item{Path: files[i].Path, Title: files[i].Title}
	}

	m.list.SetItems(items)
	return nil
}
