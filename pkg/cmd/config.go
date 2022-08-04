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
		if err := loadConfigFromPath(k, path); err != nil {
			return nil, err
		}
	} else if globalPath, err := loadGlobalConfig(k); err != nil {
		return nil, err
	} else {
		path = globalPath
	}

	var cfgFile ConfigFile
	if err := k.UnmarshalWithConf("", &cfgFile, koanf.UnmarshalConf{Tag: "json", DecoderConfig: mapStructureDecoderConfig(&cfgFile)}); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	return buildConfig(path, cfgFile), nil
}

func buildConfig(path string, cfgFile ConfigFile) *Config {
	cfgFile.Styles.apply()

	theoryCfg := theory.NewConfig(cfgFile.Theory)

	return &Config{
		ConfigFile: cfgFile,
		Theory:     theory.New(theoryCfg),

		configFilePath: path,
	}
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
			MaxColumns: 3,
			BoundaryColor: ConfigStyleColor{
				Light: "#D9DCCF",
				Dark:  "#383838",
			},
			Chord: ConfigStyle{
				Foreground: ConfigStyleColor{
					Light: "#43BF6D",
					Dark:  "#73F59F",
				},
				Italic: true,
			},
			Directive: ConfigStyle{
				Foreground: ConfigStyleColor{
					Light: "#666666",
					Dark:  "#CCCCCC",
				},
			},
			SectionName: ConfigStyle{
				Foreground: ConfigStyleColor{
					Light: "#874BFD",
					Dark:  "#7D56F4",
				},
				Underline: true,
			},
			Lyrics: ConfigStyle{
				Foreground: ConfigStyleColor{
					Light: "#111111",
					Dark:  "#EEEEEE",
				},
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
	if err := loadConfigFromPath(k, filepath.Join(configDir, "config.json")); err == nil || !errors.Is(err, os.ErrNotExist) {
		return filepath.Join(configDir, "config.json"), err
	}
	if err := loadConfigFromPath(k, filepath.Join(configDir, "config.yaml")); err == nil || !errors.Is(err, os.ErrNotExist) {
		return filepath.Join(configDir, "config.yaml"), err
	}
	if err := loadConfigFromPath(k, filepath.Join(configDir, "config.toml")); err == nil || !errors.Is(err, os.ErrNotExist) {
		return filepath.Join(configDir, "config.toml"), err
	}

	return filepath.Join(configDir, "config.toml"), nil
}

func loadConfigFromPath(k *koanf.Koanf, path string) error {
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
	MaxColumns int `json:"maxColumns,omitempty"`

	BoundaryColor ConfigStyleColor `json:"boundaryColor,omitempty"`
	Chord         ConfigStyle      `json:"chord,omitempty"`
	Directive     ConfigStyle      `json:"directive,omitempty"`
	Lyrics        ConfigStyle      `json:"lyrics,omitempty"`
	SectionName   ConfigStyle      `json:"sectionName,omitempty"`
	TitleStyle    ConfigStyle      `json:"titleStyle,omitempty"`
}

func (cs ConfigStyles) apply() {
	headerStyle = headerStyle.BorderForeground(cs.BoundaryColor.lipglossColor())
	headerFooterBoundaryStyle = headerFooterBoundaryStyle.Foreground(cs.BoundaryColor.lipglossColor())
	footerStyle = footerStyle.BorderForeground(cs.BoundaryColor.lipglossColor())

	chordStyle = cs.Chord.Apply(chordStyle)
	directiveStyle = cs.Directive.Apply(directiveStyle)
	lyricsStyle = cs.Lyrics.Apply(lyricsStyle)
	sectionNameStyle = cs.SectionName.Apply(sectionNameStyle)
	titleStyle = cs.TitleStyle.Apply(titleStyle)
}

type ConfigStyle struct {
	Background ConfigStyleColor `json:"background,omitempty"`
	Bold       bool             `json:"bold"`
	Foreground ConfigStyleColor `json:"foreground,omitempty"`
	Italic     bool             `json:"italic"`
	Underline  bool             `json:"underline"`

	f func(string) string
}

type ConfigStyleColor struct {
	Light string `json:"light,omitempty"`
	Dark  string `json:"dark,omitempty"`
}

func (csc ConfigStyleColor) lipglossColor() lipgloss.TerminalColor {
	switch {
	case len(csc.Light) != 0 && len(csc.Dark) != 0:
		return lipgloss.AdaptiveColor{
			Light: csc.Light,
			Dark:  csc.Dark,
		}
	case len(csc.Light) != 0:
		return lipgloss.Color(csc.Light)
	case len(csc.Dark) != 0:
		return lipgloss.Color(csc.Dark)
	default:
		return lipgloss.NoColor{}
	}
}

func (s *ConfigStyle) Apply(style lipgloss.Style) lipgloss.Style {
	return style.Background(s.Background.lipglossColor()).
		Bold(s.Bold).
		Foreground(s.Foreground.lipglossColor()).
		Italic(s.Italic).
		Underline(s.Underline)
}

type Config struct {
	ConfigFile
	Theory *theory.Theory

	configFilePath string
}
