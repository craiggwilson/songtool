package cmd

import (
	"fmt"

	"github.com/craiggwilson/songtool/pkg/theory"
	"github.com/jwalton/gchalk"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/rawbytes"
)

var defaultConfig = `{
	"styles": {
		"chord": {
			"color": "#69b797",
			"italic": true
		},
		"directive": {
			"color": "#666666"
		},
		"section": {
			"color": "#41829f",
			"underline": true
		},
		"text": {
			"color": "#AAAAAA"
		}
	}
}
`

func LoadConfig(path string) (*Config, error) {
	k := koanf.New(".")
	if err := k.Load(rawbytes.Provider([]byte(defaultConfig)), json.Parser()); err != nil {
		return nil, fmt.Errorf("parsing default config: %w", err)
	}

	var configFile ConfigFile
	if err := k.Unmarshal("", &configFile); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	return &Config{
		ConfigFile: configFile,
		Theory:     theory.DefaultConfig(),
	}, nil
}

type ConfigFile struct {
	Styles Styles `koanf:"styles"`
}

type Styles struct {
	Chord     Style `koanf:"chord"`
	Directive Style `koanf:"directive"`
	Section   Style `koanf:"section"`
	Text      Style `koanf:"text"`
}

type Style struct {
	Background string `koanf:"background"`
	Bold       bool   `koanf:"bold"`
	Color      string `koanf:"color"`
	Italic     bool   `koanf:"italic"`
	Underline  bool   `koanf:"underline"`

	f gchalk.ColorFn
}

func (s *Style) Render(str string) string {
	if s.f == nil {
		bldr := gchalk.New()

		if len(s.Background) > 0 {
			bldr, _ = bldr.WithStyle("bg" + s.Background)
		}

		if s.Bold {
			bldr = bldr.WithBold()
		}

		if len(s.Color) > 0 {
			bldr, _ = bldr.WithStyle(s.Color)
		}

		if s.Italic {
			bldr = bldr.WithItalic()
		}

		if s.Underline {
			bldr = bldr.WithUnderline()
		}

		s.f = bldr.StyleMust()
	}

	return s.f(str)
}

func (s *Style) Renderf(format string, a ...interface{}) string {
	return s.Render(fmt.Sprintf(format, a...))
}

type Config struct {
	ConfigFile
	Theory theory.Config
}
