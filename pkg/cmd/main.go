package cmd

import "github.com/alecthomas/kong"

var main struct {
	Key KeyCmd `cmd:"" help:"tools for keys"`
}

func Run(versionInfo VersionInfo, args []string) {
	ctx := kong.Parse(&main)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
