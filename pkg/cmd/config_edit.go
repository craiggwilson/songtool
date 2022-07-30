package cmd

import "os/exec"

type ConfigEditCmd struct{}

func (cmd *ConfigEditCmd) Run(cfg *Config) error {
	args := append([]string{cfg.configFilePath}, cfg.Edit.Args...)
	exe := exec.Command(cfg.Edit.Command, args...)
	return exe.Run()
}
