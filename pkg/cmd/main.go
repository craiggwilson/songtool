package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/craiggwilson/songtool/pkg/cmd/internal"
	"github.com/craiggwilson/songtool/pkg/cmd/internal/config"

	"github.com/alecthomas/kong"
	"github.com/muesli/termenv"
)

var mainCmd struct {
	App       internal.AppCmd       `cmd:"" help:"Loads the songtool interactive TUI." default:"withargs"`
	Cat       internal.CatCmd       `cmd:"" help:"Displays a song."`
	Chords    internal.ChordsCmd    `cmd:"" help:"Tools for working with chords."`
	Config    internal.ConfigCmd    `cmd:"" help:"Tools for managin the config."`
	Keys      internal.KeysCmd      `cmd:"" help:"Tools for working with keys."`
	Meta      internal.MetaCmd      `cmd:"" help:"Displays the meta information about a song."`
	Scales    internal.ScalesCmd    `cmd:"" help:"Tools for working with scales."`
	Transpose internal.TransposeCmd `cmd:"" help:"Transposes a song."`

	LogFile string `name:"logfile"`
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

	if len(mainCmd.LogFile) > 0 {
		f, err := os.OpenFile(mainCmd.LogFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Fprintln(parser.Stdout, fmt.Errorf("opening log file %q: %w", mainCmd.LogFile, err))
			return 4
		}
		defer f.Close()
		log.SetOutput(f)
	} else {
		log.SetOutput(io.Discard)
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
