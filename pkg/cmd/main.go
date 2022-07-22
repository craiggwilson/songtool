package cmd

import (
	"errors"
	"fmt"

	"github.com/alecthomas/kong"
)

var main struct {
	Cat       CatCmd       `cmd:"" help:"Displays a song; can be used to ensure that the song can be parsed."`
	Keys      KeysCmd      `cmd:"" help:"Lists the keys that can be used in songs."`
	Meta      MetaCmd      `cmd:"" help:"Displays the meta information about a song."`
	Transpose TransposeCmd `cmd:"" help:"Transposes a song."`
}

func Run(versionInfo VersionInfo, args []string) int {
	parser, err := kong.New(&main, kong.UsageOnError())
	if err != nil {
		panic(err)
	}

	ctx, err := parser.Parse(args)
	if err != nil {
		fmt.Fprintln(parser.Stdout, err)
		fmt.Fprintln(parser.Stdout)

		var parseErr *kong.ParseError
		if errors.As(err, &parseErr) {
			parseErr.Context.PrintUsage(false)
		}

		return 1
	}

	if err = ctx.Run(); err != nil {
		fmt.Fprintln(parser.Stdout, err)
		return 2
	}

	return 0
}
