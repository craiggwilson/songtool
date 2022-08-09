package eval

import (
	"fmt"
	"strconv"

	"github.com/alecthomas/kong"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/app/message"
	"github.com/craiggwilson/songtool/pkg/theory/interval"
	"github.com/mattn/go-shellwords"
)

func run(ctx Context, text string) tea.Cmd {
	parser, err := kong.New(&mainCmd, kong.NoDefaultHelp())
	if err != nil {
		return message.UpdateStatusError(err)
	}

	args, err := shellwords.Parse(text)
	if err != nil {
		return message.UpdateStatusError(err)
	}

	kctx, err := parser.Parse(args)
	if err != nil {
		return message.UpdateStatusError(err)
	}

	var result tea.Cmd

	err = kctx.Run(ctx, &result)
	if err != nil {
		return message.UpdateStatusError(err)
	}

	return result
}

var mainCmd struct {
	Enharmonic enharmonicCmd `cmd:"" aliases:"e" help:"Tranpose the song to it's enhmarmonic."`
	Quit       quitCmd       `cmd:"" aliases:"q" help:"Quit the app."`
	Transpose  transposeCmd  `cmd:"" aliases:"t" help:"Transpose the current song."`
}

type enharmonicCmd struct{}

func (cmd *enharmonicCmd) Run(ctx Context, result *tea.Cmd) error {
	if ctx.Meta == nil {
		return fmt.Errorf("no song loaded")
	}
	if ctx.Meta.Key == nil {
		return fmt.Errorf("current key is unset")
	}

	*result = message.TransposeSong(ctx.Meta.Key.Enharmonic())
	return nil
}

type quitCmd struct{}

func (cmd *quitCmd) Run(ctx Context, result *tea.Cmd) error {
	*result = tea.Quit
	return nil
}

type transposeCmd struct {
	Arg string `arg:"<key or step>" required:""`
}

func (cmd *transposeCmd) Run(ctx Context, result *tea.Cmd) error {
	if ctx.Meta == nil {
		return fmt.Errorf("no song loaded")
	}
	if ctx.Meta.Key == nil {
		return fmt.Errorf("current key is unset")
	}

	var intval interval.Interval
	step, err := strconv.Atoi(cmd.Arg)
	if err == nil {
		intval = ctx.Meta.Key.Step(step)
	} else {
		toKey, err := ctx.Theory.ParseKey(cmd.Arg)
		if err != nil {
			return fmt.Errorf("invalid to-key: %w", err)
		}
		intval = ctx.Meta.Key.Note().Interval(toKey.Note())
	}

	*result = message.TransposeSong(intval)
	return nil
}
