package cmd

import (
	"encoding/json"
	"fmt"
	"os"

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

	song := cmd.openSong()

	if cmd.NoChords {
		song = songio.RemoveChords(song)
	}

	if cmd.JSON {
		return cmd.printJSON(song)
	}

	return cmd.print(cfg, song)
}

func (cmd *CatCmd) print(cfg *Config, song songio.Song) error {
	_, err := WriteChordsOverLyricsWithHighlighter(&cfg.Styles, song, os.Stdout)
	return err
}

func (cmd *CatCmd) printJSON(song songio.Song) error {
	lines, err := songio.ReadAllLines(song)
	if err != nil {
		return err
	}

	out, err := json.MarshalIndent(lines, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}
