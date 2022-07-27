package cmd

import (
	"fmt"

	"github.com/craiggwilson/songtool/pkg/theory2/key"
)

type KeysCmd struct {
	JSON bool   `name:"json" help:"Prints the output as JSON."`
	Kind string `name:"kind" enum:"all,major,minor" default:"all" help:"Indicates which kind of keys to generate; major, minor, or all"`
}

func (cmd *KeysCmd) Run(cfg *Config) error {
	keys := key.List()

	if cmd.Kind == "major" {
		majorKeys := make([]key.Key, 0, len(keys)/2)
		for _, k := range keys {
			if k.Kind() == key.KindMajor {
				majorKeys = append(majorKeys, k)
			}
		}
		keys = majorKeys
	}

	if cmd.Kind == "minor" {
		minorKeys := make([]key.Key, 0, len(keys)/2)
		for _, k := range keys {
			if k.Kind() == key.KindMajor {
				minorKeys = append(minorKeys, k)
			}
		}
		keys = minorKeys
	}

	key.Sort(keys)

	if cmd.JSON {
		return cmd.printJSON(cfg, keys)
	}

	return cmd.print(cfg, keys)
}

func (cmd *KeysCmd) print(cfg *Config, keys []key.Key) error {
	for _, k := range keys {
		fmt.Println(cfg.Theory2.NameKey(k))
	}

	return nil
}

func (cmd *KeysCmd) printJSON(cfg *Config, keys []key.Key) error {
	keySurs := make([]keySurrogate, 0, len(keys))
	for _, k := range keys {
		keySurs = append(keySurs, keySurrogate{
			Name:        cfg.Theory2.NameKey(k),
			DegreeClass: k.Note().DegreeClass(),
			PitchClass:  k.Note().PitchClass(),
			Kind:        k.Kind(),
		})
	}

	return printJSON(keySurs)
}
