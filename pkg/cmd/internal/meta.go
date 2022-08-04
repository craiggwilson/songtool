package internal

import (
	"encoding/json"
	"fmt"

	"github.com/craiggwilson/songtool/pkg/cmd/internal/config"
	"github.com/craiggwilson/songtool/pkg/songio"
)

type MetaCmd struct {
	songCmd

	JSON bool `name:"json" help:"Prints the output as JSON."`
}

func (cmd *MetaCmd) Run(cfg *config.Config) error {
	defer cmd.ensurePath().Close()

	song := cmd.openSong(cfg)

	meta, err := songio.ReadMeta(cfg.Theory, song, true)
	if err != nil {
		return err
	}

	if cmd.JSON {
		return cmd.printJSON(meta)
	}

	return cmd.print(cfg, meta)
}

func (cmd *MetaCmd) print(cfg *config.Config, meta songio.Meta) error {
	if len(meta.Title) > 0 {
		fmt.Println("Title:", meta.Title)
	} else {
		fmt.Println("Title:", "<none>")
	}

	if meta.Key != nil {
		fmt.Println("Key:", meta.Key.Name)
	} else {
		fmt.Println("Key:", "<none>")
	}

	if len(meta.Sections) > 0 {
		fmt.Print("Sections: ")
		for i, section := range meta.Sections {
			if i != 0 {
				fmt.Print(", ")
			}
			fmt.Print(section)
		}
		fmt.Println()
	}

	if len(meta.Chords) > 0 {
		fmt.Print("Chords: ")
		for i, chord := range meta.Chords {
			if i != 0 {
				fmt.Print(", ")
			}
			fmt.Print(chord.Name)
		}
		fmt.Println()
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
