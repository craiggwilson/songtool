package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/craiggwilson/songtool/pkg/theory"
	"github.com/kirsle/configdir"
	"github.com/knadh/koanf"
	koanfjson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/mitchellh/mapstructure"
)

func Load(path string) (*Config, error) {
	k := koanf.New(".")
	if err := loadDefault(k); err != nil {
		return nil, fmt.Errorf("parsing default config: %w", err)
	}

	if len(path) > 0 {
		if err := loadFromPath(k, path); err != nil {
			return nil, err
		}
	} else if globalPath, err := loadGlobal(k); err != nil {
		return nil, err
	} else {
		path = globalPath
	}

	var cfgFile File
	if err := k.UnmarshalWithConf("", &cfgFile, koanf.UnmarshalConf{Tag: "json", DecoderConfig: mapStructureDecoderConfig(&cfgFile)}); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	theoryCfg := theory.NewConfig(cfgFile.Theory)

	return &Config{
		File:   cfgFile,
		Theory: theory.New(theoryCfg),

		Path: path,
	}, nil
}

func loadDefault(k *koanf.Koanf) error {
	cfg := File{
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

	marshaled, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return k.Load(rawbytes.Provider(marshaled), koanfjson.Parser())
}

func loadGlobal(k *koanf.Koanf) (string, error) {
	configDir := configdir.LocalConfig("songtool")
	if err := configdir.MakePath(configDir); err != nil {
		return "", err
	}
	if err := loadFromPath(k, filepath.Join(configDir, "config.json")); err == nil || !errors.Is(err, os.ErrNotExist) {
		return filepath.Join(configDir, "config.json"), err
	}
	if err := loadFromPath(k, filepath.Join(configDir, "config.yaml")); err == nil || !errors.Is(err, os.ErrNotExist) {
		return filepath.Join(configDir, "config.yaml"), err
	}
	if err := loadFromPath(k, filepath.Join(configDir, "config.toml")); err == nil || !errors.Is(err, os.ErrNotExist) {
		return filepath.Join(configDir, "config.toml"), err
	}

	return filepath.Join(configDir, "config.toml"), nil
}

func loadFromPath(k *koanf.Koanf, path string) error {
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
