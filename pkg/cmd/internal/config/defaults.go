package config

import (
	"os"

	"github.com/craiggwilson/songtool/pkg/theory"
)

var def = File{
	Edit: Edit{
		Command: firstNonEmptyString(
			os.Getenv("VISUAL"),
			os.Getenv("EDITOR"),
			defaultEditor,
		),
	},
	Styles: Styles{
		MaxColumns: 3,
		BoundaryColor: Color{
			Dark: "0",
		},
		Chord: Style{
			Foreground: Color{
				Dark: "8",
			},
			Italic: true,
		},
		Directive: Style{
			Foreground: Color{
				Dark: "0",
			},
		},
		SectionName: Style{
			Foreground: Color{
				Dark: "5",
			},
			Underline: true,
		},
		Lyrics: Style{},
	},
	Theory: theory.DefaultConfigBase(),
}
