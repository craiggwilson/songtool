package cmd

import (
	"os"

	"github.com/craiggwilson/songtool/pkg/songio"
)

type songCmd struct {
	Format string   `name:"format" enum:"auto,chordsOverLyrics" default:"auto" help:"The format of the song; defaults to 'auto'."`
	Path   *os.File `arg:"" optional:"" help:"The path to the song; '-' can be used for stdin."`
}

func (cmd *songCmd) ensurePath() *os.File {
	if cmd.Path == nil {
		cmd.Path = os.Stdin
	}
	return cmd.Path
}

func (cmd *songCmd) openSong(cfg *Config) songio.Song {
	return songio.ReadChordsOverLyrics(cfg.Theory, cfg.Theory, cmd.Path)
}
