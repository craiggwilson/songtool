package internal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/config"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/models"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/models/message"
)

type ViewCmd struct {
	Format string `name:"format" enum:"auto,chordsOverLyrics" default:"auto" help:"The format of the song; defaults to 'auto'."`
	Path   string `arg:"" optional:"" type:"path" help:"The path to the song."`
}

func (cmd *ViewCmd) Run(cfg *config.Config) error {
	appModel := models.NewApp(cfg, message.LoadSong(cmd.Path))

	p := tea.NewProgram(
		appModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	return p.Start()
}
