package cmd

import "fmt"

type ConfigPathCmd struct{}

func (cmd *ConfigPathCmd) Run(cfg *Config) error {
	fmt.Println(cfg.configFilePath)
	return nil
}
