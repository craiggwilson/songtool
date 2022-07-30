package cmd

import (
	"fmt"

	"github.com/craiggwilson/songtool/pkg/theory/chord"
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

func (cmd *ChordsParseCmd) print(cfg *Config, c chord.Named) error {
	fmt.Println("Name:", c.Name)
	formalName := cfg.Theory.NameChord(c.Chord)
	if formalName != c.Name {
		fmt.Println("Normalized Name:", formalName)
	}
	fmt.Println("Intervals:", c.Parsed.Chord.Intervals())
	return nil
}

func (cmd *ChordsParseCmd) printJSON(cfg *Config, c chord.Named) error {
	return printJSON(c)
}
