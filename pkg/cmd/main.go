package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/alecthomas/kong"
	"github.com/muesli/termenv"
)

var mainCmd struct {
	Cat       CatCmd       `cmd:"" help:"Displays a song."`
	Chords    ChordsCmd    `cmd:"" help:"Tools for working with chords."`
	Config    ConfigCmd    `cmd:"" help:"Tools for managin the config."`
	Keys      KeysCmd      `cmd:"" help:"Tools for working with keys."`
	Meta      MetaCmd      `cmd:"" help:"Displays the meta information about a song."`
	Scales    ScalesCmd    `cmd:"" help:"Tools for working with scales."`
	Transpose TransposeCmd `cmd:"" help:"Transposes a song."`
	View      ViewCmd      `cmd:"" help:"View a song in a modal."`
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

	cfg, err := LoadConfig("")
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
