package cmd

import (
	"github.com/craiggwilson/songtool/pkg/songio"
)

type CatCmd struct {
	songCmd

	NoChords bool `name:"no-chords" help:"Hides chords from the output."`
	JSON     bool `name:"json" xor:"json" help:"Prints the output as JSON."`
	Color    bool `name:"color" xor:"json" negatable:"" help:"Indicates whether to use color"`
}

func (cmd *CatCmd) Run(cfg *Config) error {
	defer cmd.ensurePath().Close()

	song := cmd.openSong(cfg)

	if cmd.NoChords {
		song = songio.RemoveChords(song)
	}

	if cmd.JSON {
		return cmd.printSongJSON(song)
	}

	return cmd.printSong(&cfg.Styles, song)
}
