package internal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/message"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/config"
)

type AppCmd struct {
	Format string `name:"format" enum:"auto,chordsOverLyrics" default:"auto" help:"The format of the song; defaults to 'auto'."`
	Path   string `arg:"" optional:"" type:"path" help:"The path to the song."`
}

func (cmd *AppCmd) Run(cfg *config.Config) error {
	appModel := app.New(cfg, message.LoadSong(cmd.Path))

	p := tea.NewProgram(
		appModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	return p.Start()
}
