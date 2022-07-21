package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/craiggwilson/songtool/pkg/theory"
)

type KeyCmd struct {
	List KeyListCmd `cmd:"" help:"list keys"`
}

type KeyListCmd struct {
	JSON  bool   `name:"json" xor:"style" help:"print as json"`
	Table bool   `name:"table" xor:"style" help:"print as a table"`
	Kind  string `name:"kind" enum:"major,minor,both" default:"both" help:"kind of keys; major, minor, or both"`
}

func (cmd *KeyListCmd) Run() error {
	var keys []theory.Key
	if cmd.Kind == "both" || cmd.Kind == "major" {
		keys = append(keys, theory.GenerateKeys(nil, theory.KeyMajor)...)
	}

	if cmd.Kind == "both" || cmd.Kind == "minor" {
		keys = append(keys, theory.GenerateKeys(nil, theory.KeyMinor)...)
	}

	theory.SortKeys(keys)

	if cmd.JSON {
		return cmd.printJSON(keys)
	}

	return cmd.print(keys)
}

func (cmd *KeyListCmd) print(keys []theory.Key) error {
	for _, k := range keys {
		fmt.Println(k.Name())
	}

	return nil
}

func (cmd *KeyListCmd) printJSON(keys []theory.Key) error {
	out, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}

func (cmd *KeyListCmd) printTable(keys []theory.Key) error {
	return nil
}
