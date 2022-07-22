package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/craiggwilson/songtool/pkg/songio"
)

type MetaCmd struct {
	songCmd

	JSON bool `name:"json" help:"Prints the output as JSON."`
}

func (cmd *MetaCmd) Run(cfg *Config) error {
	defer cmd.ensurePath().Close()

	song := cmd.openSong()

	meta, err := songio.ReadMeta(&cfg.Theory, song, true)
	if err != nil {
		return err
	}

	if cmd.JSON {
		return cmd.printJSON(meta)
	}

	return cmd.print(meta)
}

func (cmd *MetaCmd) print(meta songio.Meta) error {
	if len(meta.Title) > 0 {
		fmt.Println("Title:", meta.Title)
	} else {
		fmt.Println("Title:", "<none>")
	}

	if meta.Key.Note.IsValid() {
		fmt.Println("Key:", meta.Key.Name())
	} else {
		fmt.Println("Key:", "<none>")
	}

	return nil
}

func (cmd *MetaCmd) printJSON(meta songio.Meta) error {
	out, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}
