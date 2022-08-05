package models

import (
	"fmt"
	"strconv"

	"github.com/alecthomas/kong"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/craiggwilson/songtool/pkg/songio"
	"github.com/craiggwilson/songtool/pkg/theory/interval"
	"github.com/mattn/go-shellwords"
)

func runCommand(app *appModel, s string) (tea.Cmd, error) {
	parser, err := kong.New(&mainCmd, kong.NoDefaultHelp())
	if err != nil {
		return nil, err
	}

	args, err := shellwords.Parse(s)
	if err != nil {
		return nil, err
	}

	ctx, err := parser.Parse(args)
	if err != nil {
		return nil, err
	}

	var result tea.Cmd

	err = ctx.Run(app, &result)
	return result, err
}

var mainCmd struct {
	Enharmonic enharmonicCmd `cmd:"" aliases:"e" help:"Tranpose the song to it's enhmarmonic."`
	Quit       quitCmd       `cmd:"" aliases:"q" help:"Quit the app."`
	Transpose  transposeCmd  `cmd:"" aliases:"t" help:"Transpose the current song."`
}

type enharmonicCmd struct{}

func (cmd *enharmonicCmd) Run(app *appModel, _ *tea.Cmd) error {
	if app.meta.Key == nil {
		return fmt.Errorf("current key is unset")
	}

	song := songio.NewMemory(app.lines)

	intval := app.meta.Key.Enharmonic()

	transposer := songio.Transpose(app.cfg.Theory, song, intval)

	return app.SetSong(app.meta.Title, transposer)
}

type quitCmd struct{}

func (cmd *quitCmd) Run(_ *appModel, result *tea.Cmd) error {
	*result = tea.Quit
	return nil
}

type transposeCmd struct {
	Arg string `arg:"<key or step>" required:""`
}

func (cmd *transposeCmd) Run(app *appModel, _ *tea.Cmd) error {
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
