package cmd

import (
	"fmt"

	"github.com/craiggwilson/songtool/pkg/theory2/chord"
	"github.com/craiggwilson/songtool/pkg/theory2/interval"
	"github.com/craiggwilson/songtool/pkg/theory2/note"
)

type ChordsCmd struct {
	Parse ChordsParseCmd `cmd:"" help:"Parse a chord for validity and proper naming."`
}

type ChordsParseCmd struct {
	JSON bool `name:"json" help:"Prints the output as JSON."`

	Name string `arg:"<name>" help:"The name of the chord to parse."`
}

func (cmd *ChordsParseCmd) Run(cfg *Config) error {

	c, err := cfg.Theory2.ParseChord(cmd.Name)
	if err != nil {
		return err
	}

	if cmd.JSON {
		return cmd.printJSON(cfg, c)
	}

	return cmd.print(cfg, c)
}

func (cmd *ChordsParseCmd) print(cfg *Config, c chord.Parsed) error {
	fmt.Println(c.Name(cfg.Theory2))
	return nil
}

func (cmd *ChordsParseCmd) printJSON(cfg *Config, c chord.Parsed) error {
	return printJSON(struct {
		Name              string              `json:"name"`
		Root              note.Note           `json:"root"`
		Suffix            string              `json:"suffix"`
		BaseNoteDelimiter string              `json:"baseNoteDelimiter,omitempty"`
		Base              *note.Note          `json:"base"`
		Intervals         []interval.Interval `json:"intervals"`
	}{c.Name(cfg.Theory2), c.Root(), c.Suffix, c.BaseNoteDelimiter, c.Base(), c.Intervals()})
}
