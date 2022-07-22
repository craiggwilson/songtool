package cmd

import (
	"os"

	"github.com/craiggwilson/songtool/pkg/songio"
)

type CatCmd struct {
	songCmd
}

func (cmd *CatCmd) Run() error {
	defer cmd.ensurePath().Close()

	song := cmd.openSong()

	_, err := songio.WriteChordsOverLyrics(nil, song, os.Stdout)
	return err
}
