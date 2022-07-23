package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/craiggwilson/songtool/pkg/theory"
)

type KeysCmd struct {
	JSON bool   `name:"json" help:"Prints the output as JSON."`
	Kind string `name:"kind" enum:"all,major,minor" default:"all" help:"Indicates which kind of keys to generate; major, minor, or all"`
}

func (cmd *KeysCmd) Run(cfg *Config) error {
	var keys []theory.Key
	if cmd.Kind == "all" || cmd.Kind == "major" {
		keys = append(keys, cfg.Theory.GenerateKeys(theory.KeyMajor)...)
	}

	if cmd.Kind == "all" || cmd.Kind == "minor" {
		keys = append(keys, cfg.Theory.GenerateKeys(theory.KeyMinor)...)
	}

	theory.SortKeys(keys)

	if cmd.JSON {
		return cmd.printJSON(keys)
	}

	return cmd.print(keys)
}

func (cmd *KeysCmd) print(keys []theory.Key) error {
	for _, k := range keys {
		fmt.Println(k.Name())
	}

	return nil
}

func (cmd *KeysCmd) printJSON(keys []theory.Key) error {
	out, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}
