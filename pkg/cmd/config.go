package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jwalton/gchalk"
	"github.com/kirsle/configdir"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"

	"github.com/craiggwilson/songtool/pkg/theory"
)

type ConfigCmd struct {
	Cat ConfigCatCmd `cmd:"" help:"Prints the config."`
}

func LoadConfig(path string) (*Config, error) {
	k := koanf.New(".")
	if err := loadDefaultConfig(k); err != nil {
		return nil, fmt.Errorf("parsing default config: %w", err)
	}

	if len(path) > 0 {
		if err := loadConfig(k, path); err != nil {
			return nil, err
		}
	} else if err := loadGlobalConfig(k); err != nil {
		return nil, err
	}

	var configFile ConfigFile
	if err := k.Unmarshal("", &configFile); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	theoryCfg := theory.NewConfig(configFile.Theory)

	return &Config{
		ConfigFile: configFile,
		Theory:     theory.New(theoryCfg),
	}, nil
}

func loadDefaultConfig(k *koanf.Koanf) error {
	// var defaultConfigFile = `
	// [styles.chord]
	// color = "#69b797"
	// italic = true

	// [styles.directive]
	// color = "#666666"

	// [styles.section]
	// color = "#41829f"
	// underline = true

	// [styles.text]
	// color = "#AAAAAA"
	// `
	//k.Load(rawbytes.Provider([]byte(defaultConfigFile)), toml.Parser());

	cfg := ConfigFile{
		Styles: Styles{
			Chord: Style{
				Color:  "#69b797",
				Italic: true,
			},
			Directive: Style{
				Color: "#666666",
			},
			Section: Style{
				Color:     "#41829f",
				Underline: true,
			},
			Text: Style{
				Color: "#AAAAAA",
			},
		},
		Theory: theory.DefaultConfigBase(),
	}

	return k.Load(structs.Provider(cfg, ""), nil)
}

func loadGlobalConfig(k *koanf.Koanf) error {
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
	Styles Styles            `json:"styles"`
	Theory theory.ConfigBase `json:"theory"`
}

type Styles struct {
	Chord     Style `json:"chord,omitempty"`
	Directive Style `json:"directive,omitempty"`
	Section   Style `json:"section,omitempty"`
	Text      Style `json:"text,omitempty"`
}

type Style struct {
	Background string `json:"background,omitempty"`
	Bold       bool   `json:"bold"`
	Color      string `json:"color,omitempty"`
	Italic     bool   `json:"italic"`
	Underline  bool   `json:"underline"`

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
	Theory *theory.Theory
}
