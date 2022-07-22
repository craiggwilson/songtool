package cmd

import (
	"errors"
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/fatih/color"
	"github.com/jwalton/go-supportscolor"
)

var main struct {
	Cat       CatCmd       `cmd:"" help:"Displays a song."`
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

	color.NoColor = !supportscolor.Stdout().SupportsColor
	cfg := LoadConfig("")

	if err = ctx.Run(cfg); err != nil {
		fmt.Fprintln(parser.Stdout, err)
		return 2
	}

	return 0
}
