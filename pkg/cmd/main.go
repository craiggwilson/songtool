package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/craiggwilson/songtool/pkg/cmd/internal"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/config"

	"github.com/alecthomas/kong"
	"github.com/muesli/termenv"
)

var mainCmd struct {
	Cat       internal.CatCmd       `cmd:"" help:"Displays a song."`
	Chords    internal.ChordsCmd    `cmd:"" help:"Tools for working with chords."`
	Config    internal.ConfigCmd    `cmd:"" help:"Tools for managin the config."`
	Keys      internal.KeysCmd      `cmd:"" help:"Tools for working with keys."`
	Meta      internal.MetaCmd      `cmd:"" help:"Displays the meta information about a song."`
	Scales    internal.ScalesCmd    `cmd:"" help:"Tools for working with scales."`
	Transpose internal.TransposeCmd `cmd:"" help:"Transposes a song."`
	View      internal.ViewCmd      `cmd:"" help:"View a song in a modal."`
}

func Run(versionInfo VersionInfo, args []string) int {
	parser, err := kong.New(&mainCmd, kong.UsageOnError(), kong.Vars{
		"color": strconv.FormatBool(termenv.EnvColorProfile() != termenv.Ascii),
	})
	if err != nil {
		panic(err)
	}

	ctx, err := parser.Parse(args)
	if err != nil {
		fmt.Fprintln(parser.Stdout, err)

		var parseErr *kong.ParseError
		if errors.As(err, &parseErr) {
			parseErr.Context.PrintUsage(false)
		}

		return 1
	}

	cfg, err := config.Load("")
	if err != nil {
		fmt.Fprintln(parser.Stdout, err)
		return 2
	}

	if err = ctx.Run(cfg); err != nil {
		fmt.Fprintln(parser.Stdout, err)
		return 3
	}

	return 0
}
