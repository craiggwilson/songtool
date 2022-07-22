package cmd

import (
	"fmt"
	"os"

	"github.com/craiggwilson/songtool/pkg/songio"
	"github.com/craiggwilson/songtool/pkg/theory"
)

type SongCmd struct {
	Cat       SongCatCmd       `cmd:"" help:"display a song"`
	Transpose SongTransposeCmd `cmd:"" help:"transpose a song"`
}

type SongCatCmd struct {
	Format string `name:"format" enum:"chordsOverLyrics,chordpro" default:"chordsOverLyrics"`
	Path   string `arg:"" help:"path to the song" type:"path" required:""`
}

func (cmd *SongCatCmd) Run() error {
	f, err := os.Open(cmd.Path)
	if err != nil {
		return fmt.Errorf("opening %q: %w", cmd.Path, err)
	}

	defer f.Close()

	it := songio.ReadChordsOverLyrics(nil, f)

	_, err = songio.WriteChordsOverLyrics(nil, it, os.Stdout)
	return err
}

type SongTransposeCmd struct {
	Path     string `arg:"" help:"path to the song" type:"path" required:""`
	Interval int    `name:"interval" short:"i"`
}

func (cmd *SongTransposeCmd) Run() error {
	f, err := os.Open(cmd.Path)
	if err != nil {
		return fmt.Errorf("opening %q: %w", cmd.Path, err)
	}

	defer f.Close()

	key, _ := theory.ParseKey(nil, "G")

	var it songio.LineIter = songio.ReadChordsOverLyrics(nil, f)

	it = songio.Transpose(nil, it, theory.IntervalFromStep(nil, key.Note, cmd.Interval, theory.EnharmonicSharp))

	_, err = songio.WriteChordsOverLyrics(nil, it, os.Stdout)
	return err
}
