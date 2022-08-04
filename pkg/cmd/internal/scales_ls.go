package internal

import (
	"fmt"

	"github.com/craiggwilson/songtool/pkg/cmd/internal/config"
	"github.com/craiggwilson/songtool/pkg/theory"
)

type ScalesLsCmd struct {
	JSON bool `name:"json" xor:"json" help:"Prints the output as JSON."`
}

func (cmd *ScalesLsCmd) Run(cfg *config.Config) error {
	scales := cfg.Theory.ListScales()
	theory.SortScaleMetas(scales)

	if cmd.JSON {
		return cmd.printJSON(scales)
	}

	return cmd.print(scales)
}

func (cmd *ScalesLsCmd) print(scales []theory.ScaleMeta) error {
	for _, scale := range scales {
		fmt.Println(scale.Name, scale.Intervals)
	}

	return nil
}

func (cmd *ScalesLsCmd) printJSON(scales []theory.ScaleMeta) error {
	return printJSON(scales)
}
