package internal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/config"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/models"
	"github.com/craiggwilson/songtool/pkg/songio"
)

type ViewCmd struct {
	songCmd
}

func (cmd *ViewCmd) Run(cfg *config.Config) error {
	path := cmd.ensurePath()

	song := cmd.openSong(cfg)

	lines, err := songio.ReadAllLines(song)
	if err != nil {
		return err
	}

	path.Close()

	appModel := models.NewApp(cfg, models.UpdateSongFromSource(cfg.Theory, path.Name(), songio.FromLines(lines)))

	p := tea.NewProgram(
		appModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	return p.Start()
}
