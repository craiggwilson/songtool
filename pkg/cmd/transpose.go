package cmd

import (
	"fmt"
	"os"

	"github.com/craiggwilson/songtool/pkg/songio"
	"github.com/craiggwilson/songtool/pkg/theory"
)

type TransposeCmd struct {
	songCmd

	FromKey  string `name:"from-key" help:"The current key of the song; will be discovered automatically when not specified."`
	Interval int    `name:"interval" short:"i" xor:"keyinterval" required:"" help:"The number of steps to transpose the song; can be negative. Cannot be used to 'to-key'."`
	ToKey    string `name:"to-key" xor:"keyinterval" required:"" help:"The desired key of the song. Cannot be used with 'interval'."`
}

func (cmd *TransposeCmd) Run() error {
	defer cmd.ensurePath().Close()

	song := cmd.openSong()

	var fromKey theory.Key
	if len(cmd.FromKey) == 0 {
		rewinder := songio.NewRewinder(song)
		meta, err := songio.ReadMeta(nil, rewinder, false)
		if err != nil {
			return err
		}

		fromKey = meta.Key
		if !fromKey.Note.IsValid() {
			return fmt.Errorf("could not infer from-key")
		}

		song = rewinder.Rewind()
	}

	if !fromKey.Note.IsValid() {
		fk, err := theory.ParseKey(nil, cmd.FromKey)
		if err != nil {
			return fmt.Errorf("invalid from-key: %w", err)
		}

		fromKey = fk
	}

	var interval theory.Interval
	if len(cmd.ToKey) > 0 {
		toKey, err := theory.ParseKey(nil, cmd.ToKey)
		if err != nil {
			return fmt.Errorf("invalid to-key: %w", err)
		}
		interval = theory.IntervalFromDiff(fromKey.Note, toKey.Note)
	} else {
		interval = theory.IntervalFromStep(nil, fromKey.Note, cmd.Interval, theory.EnharmonicSharp)
	}

	transposer := songio.Transpose(nil, song, interval)

	_, err := songio.WriteChordsOverLyrics(nil, transposer, os.Stdout)
	return err
}
