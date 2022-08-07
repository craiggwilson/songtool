package song

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/footer"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/header"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/songtext"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/config"
)

func New(cfg *config.Config) Model {
	header := header.New()
	header.BorderColor = cfg.Styles.BoundaryColor.Color()
	header.KeyStyle = cfg.Styles.Chord.Style()
	header.TitleStyle = cfg.Styles.Title.Style()

	songtext := songtext.New()
	songtext.MaxColumns = cfg.Styles.MaxColumns
	songtext.ChordStyle = cfg.Styles.Chord.Style()
	songtext.LyricsStyle = cfg.Styles.Lyrics.Style()
	songtext.SectionNameStyle = cfg.Styles.SectionName.Style()

	footer := footer.New()
	footer.BorderColor = cfg.Styles.BoundaryColor.Color()
	footer.ScrollPercentStyle = cfg.Styles.Title.Style()

	return Model{
		header:   header,
		songtext: songtext,
		footer:   footer,
	}
}

type Model struct {
	Height int
	Width  int
	KeyMap KeyMap

	header   header.Model
	songtext songtext.Model
	footer   footer.Model
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.songtext.KeyMap = m.KeyMap
	m.header.Width = m.Width
	m.songtext.Width = m.Width
	m.footer.Width = m.Width
	m.footer.ScrollPercent = m.songtext.ScrollPercent()

	headerHeight := lipgloss.Height(m.header.View())
	footerHeight := lipgloss.Height(m.footer.View())
	m.songtext.Height = m.Height - headerHeight - footerHeight

	m.header, cmd = m.header.Update(msg)
	cmds = append(cmds, cmd)

	m.songtext, cmd = m.songtext.Update(msg)
	cmds = append(cmds, cmd)

	m.footer, cmd = m.footer.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.header.View(), m.songtext.View(), m.footer.View())
}
