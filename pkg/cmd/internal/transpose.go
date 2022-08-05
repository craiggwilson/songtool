package internal

import (
	"fmt"
	"os"

	"github.com/craiggwilson/songtool/pkg/cmd/internal/config"
	"github.com/craiggwilson/songtool/pkg/songio"
	"github.com/craiggwilson/songtool/pkg/theory/interval"
	"github.com/craiggwilson/songtool/pkg/theory/key"
)

type TransposeCmd struct {
	songCmd

	FromKey  string `name:"from-key" help:"The current key of the song; will be discovered automatically when not specified."`
	Interval int    `name:"interval" short:"i" xor:"keyinterval" required:"" help:"The number of steps to transpose the song; can be negative. Cannot be used to 'to-key'."`
	ToKey    string `name:"to-key" xor:"keyinterval" required:"" help:"The desired key of the song. Cannot be used with 'interval'."`

	JSON  bool `name:"json" xor:"json" help:"Prints the output as JSON."`
	Color bool `name:"color" xor:"json" negatable:"" help:"Indicates whether to use color"`
}

func (cmd *TransposeCmd) Run(cfg *config.Config) error {
	defer cmd.ensurePath().Close()

	song := cmd.openSong(cfg)

	var fromKey *key.Named
	if len(cmd.FromKey) == 0 {
		rewinder := songio.NewRewinder(song)
		meta, err := songio.ReadMeta(cfg.Theory, rewinder, false)
		if err != nil {
			return err
		}

		fromKey = meta.Key
		if fromKey == nil {
			return fmt.Errorf("could not infer from-key")
		}

		song = rewinder.Rewind()
	}

	if fromKey == nil {
		fk, err := cfg.Theory.ParseKey(cmd.FromKey)
		if err != nil {
			return fmt.Errorf("invalid from-key: %w", err)
		}

		fromKey = &fk
	}

	var intval interval.Interval
	if len(cmd.ToKey) > 0 {
		toKey, err := cfg.Theory.ParseKey(cmd.ToKey)
		if err != nil {
			return fmt.Errorf("invalid to-key: %w", err)
		}
		intval = fromKey.Note().Interval(toKey.Note())
	} else {
		intval = interval.FromStep(cmd.Interval)
	}

	transposer := songio.Transpose(cfg.Theory, song, intval)

	_, err := songio.WriteChordsOverLyrics(cfg.Theory, transposer, os.Stdout)
	return err
}
