package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/craiggwilson/songtool/pkg/theory"
	"github.com/craiggwilson/songtool/pkg/theory2"
	"github.com/jwalton/gchalk"
	"github.com/kirsle/configdir"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
)

var defaultConfigFile = `
[styles.chord]
color = "#69b797"
italic = true

[styles.directive]
color = "#666666"

[styles.section]
color = "#41829f"
underline = true

[styles.text]
color = "#AAAAAA"
`

func LoadConfig(path string) (*Config, error) {
	k := koanf.New(".")
	if err := k.Load(rawbytes.Provider([]byte(defaultConfigFile)), toml.Parser()); err != nil {
		return nil, fmt.Errorf("parsing default config: %w", err)
	}

	if len(path) > 0 {
		if err := loadConfig(k, path); err != nil {
			return nil, err
		}
	} else if err := loadDefaultConfig(k); err != nil {
		return nil, err
	}

	var configFile ConfigFile
	if err := k.Unmarshal("", &configFile); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	return &Config{
		ConfigFile: configFile,
		Theory:     theory.NewDefault(),
		Theory2:    theory2.Default(),
	}, nil
}

func loadDefaultConfig(k *koanf.Koanf) error {
	configDir := configdir.LocalConfig("songtool")
	if err := configdir.MakePath(configDir); err != nil {
		return err
	}
	if err := loadConfig(k, filepath.Join(configDir, "config.json")); err == nil || !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := loadConfig(k, filepath.Join(configDir, "config.yaml")); err == nil || !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := loadConfig(k, filepath.Join(configDir, "config.toml")); err == nil || !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

func loadConfig(k *koanf.Koanf, path string) error {
	ext := filepath.Ext(path)
	switch ext {
	case ".json":
		if err := k.Load(file.Provider(path), json.Parser()); err != nil {
			return fmt.Errorf("loading config at %q: %w", path, err)
		}
	case ".yaml":
		if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
			return fmt.Errorf("loading config at %q: %w", path, err)
		}
	case ".toml":
		if err := k.Load(file.Provider(path), toml.Parser()); err != nil {
			return fmt.Errorf("loading config at %q: %w", path, err)
		}
	default:
		return fmt.Errorf("unsupported config format %q", ext)
	}

	return nil
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
	Theory  *theory.Theory
	Theory2 *theory2.Theory
}
