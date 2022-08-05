package models

import (
	"fmt"
	"strconv"

	"github.com/alecthomas/kong"
	"github.com/craiggwilson/songtool/pkg/songio"
	"github.com/craiggwilson/songtool/pkg/theory/interval"
	"github.com/mattn/go-shellwords"
)

func runCommand(app *appModel, s string) error {
	parser, err := kong.New(&mainCmd, kong.NoDefaultHelp())
	if err != nil {
		return err
	}

	args, err := shellwords.Parse(s)
	if err != nil {
		return err
	}

	ctx, err := parser.Parse(args)
	if err != nil {
		return err
	}

	return ctx.Run(app)
}

var mainCmd struct {
	Enharmonic enharmonicCmd `cmd:"" aliases:"e" help:"Tranpose the song to it's enhmarmonic."`
	Transpose  transposeCmd  `cmd:"" aliases:"t" help:"Transpose the current song."`
}

type enharmonicCmd struct{}

func (cmd *enharmonicCmd) Run(app *appModel) error {
	if app.meta.Key == nil {
		return fmt.Errorf("current key is unset")
	}

	song := songio.NewMemory(app.lines)

	intval := app.meta.Key.Enharmonic()

	transposer := songio.Transpose(app.cfg.Theory, song, intval)

	return app.SetSong(app.meta.Title, transposer)
}

type transposeCmd struct {
	Arg string `arg:"<key or step>" required:""`
}

func (cmd *transposeCmd) Run(app *appModel) error {
	if app.meta.Key == nil {
		return fmt.Errorf("current key is unset")
	}

	song := songio.NewMemory(app.lines)

	var intval interval.Interval
	step, err := strconv.Atoi(cmd.Arg)
	if err == nil {
		intval = app.meta.Key.Step(step)
	} else {
		toKey, err := app.cfg.Theory.ParseKey(cmd.Arg)
		if err != nil {
			return fmt.Errorf("invalid to-key: %w", err)
		}
		intval = app.meta.Key.Note().Interval(toKey.Note())
	}

	transposer := songio.Transpose(app.cfg.Theory, song, intval)

	return app.SetSong(app.meta.Title, transposer)
}