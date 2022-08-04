package internal

import "github.com/craiggwilson/songtool/pkg/cmd/internal/config"

type ConfigCatCmd struct {
}

func (cmd *ConfigCatCmd) Run(cfg *config.Config) error {
	return printJSON(cfg.File)
}
