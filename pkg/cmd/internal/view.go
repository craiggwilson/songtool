package internal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/config"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/models"
)

type ViewCmd struct {
	songCmd
}

func (cmd *ViewCmd) Run(cfg *config.Config) error {
	path := cmd.ensurePath()

	song := cmd.openSong(cfg)

	appModel := models.NewApp(cfg)

	appModel.SetSong(path.Name(), song)

	path.Close()

	p := tea.NewProgram(
		appModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	return p.Start()
}
