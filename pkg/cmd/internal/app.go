package internal

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/message"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/config"
)

type AppCmd struct {
	Format string `name:"format" enum:"auto,chordsOverLyrics" default:"auto" help:"The format of the song; defaults to 'auto'."`
	Path   string `arg:"" required:"" type:"path" help:"The path to the song."`
}

func (cmd *AppCmd) Run(cfg *config.Config) error {

	var appCmds []tea.Cmd
	if len(cmd.Path) > 0 {
		fi, err := os.Stat(cmd.Path)
		if err != nil {
			return fmt.Errorf("could not stat %s: %w", cmd.Path, err)
		}

		if !fi.IsDir() {
			appCmds = append(appCmds,
				message.LoadSong(cmd.Path),
				message.LoadDirectory(filepath.Dir(cmd.Path)),
				message.EnterSongMode(),
			)
		} else {
			appCmds = append(appCmds,
				message.LoadDirectory(cmd.Path),
				message.EnterExplorerMode(),
			)
		}
	}

	appModel := app.New(cfg, appCmds...)

	p := tea.NewProgram(
		appModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	return p.Start()
}
