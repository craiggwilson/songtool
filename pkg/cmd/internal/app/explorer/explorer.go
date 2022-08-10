package explorer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/message"
)

func New() Model {
	return Model{
		KeyMap: DefaultKeyMap(),
		Styles: DefaultStyles(),
	}
}

type Model struct {
	KeyMap *KeyMap
	Styles Styles
	Height int
	Width  int

	leftColumnIdx   int
	selectedItemIdx int
	items           []item
	itemsPerColumn  int

	renderedColumns []string
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
		case key.Matches(tmsg, m.KeyMap.MoveRight):
			if m.selectedItemIdx+m.itemsPerColumn < len(m.items) {
				m.updateSelectedIndex(m.selectedItemIdx + m.itemsPerColumn)
			}
		case key.Matches(tmsg, m.KeyMap.MoveLeft):
			if m.selectedItemIdx-m.itemsPerColumn >= 0 {
				m.updateSelectedIndex(m.selectedItemIdx - m.itemsPerColumn)
			}
		case key.Matches(tmsg, m.KeyMap.MoveUp):
			if m.selectedItemIdx > 0 {
				m.updateSelectedIndex(m.selectedItemIdx - 1)
			}
		case key.Matches(tmsg, m.KeyMap.MoveDown):
			if m.selectedItemIdx+1 < len(m.items) {
				m.updateSelectedIndex(m.selectedItemIdx + 1)
			}
		case key.Matches(tmsg, m.KeyMap.Select):
			if len(m.items) > 0 {
				item := m.items[m.selectedItemIdx]
				return m, tea.Batch(
					message.LoadSong(item.Path),
					message.EnterSongMode(),
				)
			}
		}
	case message.InvalidateMsg:
		m.itemsPerColumn = m.Height - 1
		m.renderColumns()
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var columnsToRender []string

	currentWidth := 0
	for i := m.leftColumnIdx; i < len(m.renderedColumns); i++ {
		width := lipgloss.Width(m.renderedColumns[i])
		if currentWidth+width > m.Width {
			break
		}
		columnsToRender = append(columnsToRender, m.renderedColumns[i])
		currentWidth += width
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, columnsToRender...)
}

func (m *Model) renderColumns() {
	m.renderedColumns = m.renderedColumns[:0]
	colCount := len(m.items)/(m.Height-1) + 1

	for i := 0; i < colCount; i++ {
		m.renderedColumns = append(m.renderedColumns, m.renderColumn(i))
	}
}

func (m *Model) renderColumn(columnIdx int) string {
	var render strings.Builder
	for i := (m.Height - 1) * columnIdx; i < len(m.items) && i < (m.Height-1)*(columnIdx+1); i++ {
		fn := m.Styles.ItemStyle.Render
		if i == m.selectedItemIdx {
			fn = m.Styles.SelectedItemStyle.Render
		}

		render.WriteString(fn(m.items[i].Text) + "\n")
	}

	return m.Styles.ColumnStyle.Render(render.String())
}

func (m *Model) updateItems(files []message.FileItem) tea.Cmd {
	items := make([]item, len(files))
	for i := range files {
		items[i].Path = files[i].Path
		title := filepath.Base(items[i].Path)
		if ext := filepath.Ext(title); len(ext) > 0 {
			title = strings.TrimSuffix(title, ext)
		}
		key := ""
		if files[i].Meta != nil {
			if len(files[i].Meta.Title) > 0 {
				title = files[i].Meta.Title
			}

			if files[i].Meta.Key != nil {
				key = files[i].Meta.Key.Name
			}
		}

		if key != "" {
			key = fmt.Sprintf("[%s]", m.Styles.KeyStyle.Render(key))
		}

		items[i].Text = lipgloss.JoinHorizontal(lipgloss.Top, fmt.Sprintf("%-5s", key), title)
	}

	m.items = items
	return message.Invalidate()
}

func (m *Model) updateSelectedIndex(to int) {
	fromColIdx := m.selectedItemIdx / m.itemsPerColumn
	toColIdx := to / m.itemsPerColumn
	m.selectedItemIdx = to

	m.renderedColumns[fromColIdx] = m.renderColumn(fromColIdx)
	if fromColIdx != toColIdx {
		m.renderedColumns[toColIdx] = m.renderColumn(toColIdx)
	}

	if toColIdx > 1 {
		m.leftColumnIdx = toColIdx - 1
	} else {
		m.leftColumnIdx = 0
	}
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
