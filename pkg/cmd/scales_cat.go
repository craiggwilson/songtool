package cmd

import (
	"fmt"

	"github.com/craiggwilson/songtool/pkg/theory/interval"
	"github.com/craiggwilson/songtool/pkg/theory/scale"
)

type ScalesCatCmd struct {
	Root string `arg:"<root>" help:"The name of the root note."`
	Name string `arg:"<name>" help:"The name of scale to view."`

	JSON bool `name:"json" help:"Prints the output as JSON."`
}

func (cmd *ScalesCatCmd) Run(cfg *Config) error {
	meta, ok := cfg.Theory.LookupScale(cmd.Name)
	if !ok {
		return fmt.Errorf("unknown scale %q", cmd.Name)
	}

	root, err := cfg.Theory.ParseNote(cmd.Root)
	if err != nil {
		return fmt.Errorf("parsing note: %w", err)
	}

	scale := scale.Generate(cfg.Theory.NameNote(root)+" "+meta.Name, root, meta.Intervals...)

	if cmd.JSON {
		return cmd.printJSON(cfg, scale)
	}

	return cmd.print(cfg, scale)
}

func (cmd *ScalesCatCmd) print(cfg *Config, scale scale.Scale) error {
	notes := scale.Notes()
	noteNames := make([]string, 0, len(notes))
	for _, n := range notes {
		noteNames = append(noteNames, cfg.Theory.NameNote(n))
	}

	fmt.Println(noteNames)
	return nil
}

func (cmd *ScalesCatCmd) printJSON(cfg *Config, scale scale.Scale) error {
	notes := scale.Notes()
	noteSurs := make([]noteSurrogate, 0, len(notes))
	pitchClassOffset := notes[0].PitchClass()
	for _, n := range notes {
		noteSurs = append(noteSurs, noteSurrogate{
			Name:        cfg.Theory.NameNote(n),
			Interval:    interval.FromStep(n.PitchClass() - pitchClassOffset).String(),
			DegreeClass: n.DegreeClass(),
			PitchClass:  n.PitchClass(),
		})
	}

	return printJSON(scaleSurrogate{
		Name:  scale.Name(),
		Notes: noteSurs,
	})
}
