package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/craiggwilson/songtool/pkg/songio"
)

type CatCmd struct {
	songCmd

	JSON bool `name:"json" help:"Prints the output as JSON."`
}

func (cmd *CatCmd) Run() error {
	defer cmd.ensurePath().Close()

	song := cmd.openSong()

	if cmd.JSON {
		return cmd.printJSON(song)
	}

	return cmd.print(song)
}

func (cmd *CatCmd) print(song songio.Song) error {
	_, err := songio.WriteChordsOverLyrics(nil, song, os.Stdout)
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
