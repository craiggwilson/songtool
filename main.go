package main

import (
	"os"

	"github.com/craiggwilson/songtool/pkg/cmd"
)

var (
	version string
	commit  string
	date    string
	builtBy string
)

func main() {
	cmd.Run(cmd.VersionInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
		BuiltBy: builtBy,
	}, os.Args[1:])
}
