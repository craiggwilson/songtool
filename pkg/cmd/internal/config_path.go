package internal

import (
	"fmt"

	"github.com/craiggwilson/songtool/pkg/cmd/internal/config"
)

type ConfigPathCmd struct{}

func (cmd *ConfigPathCmd) Run(cfg *config.Config) error {
	fmt.Println(cfg.Path)
	return nil
}
