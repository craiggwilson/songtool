package cmd

import (
	"fmt"

	"github.com/craiggwilson/songtool/pkg/theory/chord"
	"github.com/craiggwilson/songtool/pkg/theory/interval"
	"github.com/craiggwilson/songtool/pkg/theory/note"
)

type ChordsParseCmd struct {
	JSON bool `name:"json" help:"Prints the output as JSON."`

	Name string `arg:"<name>" help:"The name of the chord to parse."`
}

func (cmd *ChordsParseCmd) Run(cfg *Config) error {

	c, err := cfg.Theory.ParseChord(cmd.Name)
	if err != nil {
		return err
	}

	if cmd.JSON {
		return cmd.printJSON(cfg, c)
	}

	return cmd.print(cfg, c)
}

func (cmd *ChordsParseCmd) print(cfg *Config, c chord.Parsed) error {
	fmt.Println(c.Name(cfg.Theory))
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
	}{c.Name(cfg.Theory), c.Root(), c.Suffix, c.BaseNoteDelimiter, c.Base(), c.Intervals()})
}
