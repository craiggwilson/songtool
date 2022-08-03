package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/kirsle/configdir"
	"github.com/knadh/koanf"
	koanfjson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/mitchellh/mapstructure"

	"github.com/craiggwilson/songtool/pkg/theory"
)

type ConfigCmd struct {
	Cat  ConfigCatCmd  `cmd:"" help:"Prints the config."`
	Edit ConfigEditCmd `cmd:"" help:"Launches an editor for the config file."`
	Path ConfigPathCmd `cmd:"" help:"Prints the location of the config file."`
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
	} else if globalPath, err := loadGlobalConfig(k); err != nil {
		return nil, err
	} else {
		path = globalPath
	}

	var configFile ConfigFile
	if err := k.UnmarshalWithConf("", &configFile, koanf.UnmarshalConf{Tag: "json", DecoderConfig: mapStructureDecoderConfig(&configFile)}); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	theoryCfg := theory.NewConfig(configFile.Theory)

	return &Config{
		ConfigFile: configFile,
		Theory:     theory.New(theoryCfg),

		configFilePath: path,
	}, nil
}

func loadDefaultConfig(k *koanf.Koanf) error {
	cfg := ConfigFile{
		Edit: ConfigEdit{
			Command: firstNonEmptyString(
				os.Getenv("VISUAL"),
				os.Getenv("EDITOR"),
				defaultEditor,
			),
		},
		Styles: ConfigStyles{
			Chord: ConfigStyle{
				Color:  "#69b797",
				Italic: true,
			},
			Directive: ConfigStyle{
				Color: "#666666",
			},
			Section: ConfigStyle{
				Color:     "#41829f",
				Underline: true,
			},
			Text: ConfigStyle{
				Color: "#AAAAAA",
			},
		},
		Theory: theory.DefaultConfigBase(),
	}

	marshaled, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return k.Load(rawbytes.Provider(marshaled), koanfjson.Parser())
}

func loadGlobalConfig(k *koanf.Koanf) (string, error) {
	configDir := configdir.LocalConfig("songtool")
	if err := configdir.MakePath(configDir); err != nil {
		return "", err
	}
	if err := loadConfig(k, filepath.Join(configDir, "config.json")); err == nil || !errors.Is(err, os.ErrNotExist) {
		return filepath.Join(configDir, "config.json"), err
	}
	if err := loadConfig(k, filepath.Join(configDir, "config.yaml")); err == nil || !errors.Is(err, os.ErrNotExist) {
		return filepath.Join(configDir, "config.yaml"), err
	}
	if err := loadConfig(k, filepath.Join(configDir, "config.toml")); err == nil || !errors.Is(err, os.ErrNotExist) {
		return filepath.Join(configDir, "config.toml"), err
	}

	return filepath.Join(configDir, "config.toml"), nil
}

func loadConfig(k *koanf.Koanf, path string) error {
	ext := filepath.Ext(path)
	switch ext {
	case ".json":
		if err := k.Load(file.Provider(path), koanfjson.Parser()); err != nil {
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

func mapStructureDecoderConfig(o interface{}) *mapstructure.DecoderConfig {
	return &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.TextUnmarshallerHookFunc(),
		),
		Result:  o,
		TagName: "json",
	}
}

type ConfigFile struct {
	Edit   ConfigEdit        `json:"edit"`
	Styles ConfigStyles      `json:"styles"`
	Theory theory.ConfigBase `json:"theory"`
}

type ConfigEdit struct {
	Command string   `json:"command"`
	Args    []string `json:"args,omitempty"`
}

type ConfigStyles struct {
	Chord     ConfigStyle `json:"chord,omitempty"`
	Directive ConfigStyle `json:"directive,omitempty"`
	Section   ConfigStyle `json:"section,omitempty"`
	Text      ConfigStyle `json:"text,omitempty"`
}

type ConfigStyle struct {
	Background string `json:"background,omitempty"`
	Bold       bool   `json:"bold"`
	Color      string `json:"color,omitempty"`
	Italic     bool   `json:"italic"`
	Underline  bool   `json:"underline"`

	f func(string) string
}

func (s *ConfigStyle) Render(str string) string {
	if s.f == nil {
		bldr := lipgloss.NewStyle()

		if len(s.Background) > 0 {
			bldr = bldr.Background(lipgloss.Color(s.Background))
		}

		if s.Bold {
			bldr = bldr.Bold(true)
		}

		if len(s.Color) > 0 {
			bldr = bldr.Foreground(lipgloss.Color(s.Color))
		}

		if s.Italic {
			bldr = bldr.Italic(true)
		}

		if s.Underline {
			bldr = bldr.Underline(true)
		}

		s.f = bldr.Render
	}

	return s.f(str)
}

func (s *ConfigStyle) Renderf(format string, a ...interface{}) string {
	return s.Render(fmt.Sprintf(format, a...))
}

type Config struct {
	ConfigFile
	Theory *theory.Theory

	configFilePath string
}
