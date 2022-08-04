package internal

import (
	"os/exec"

	"github.com/craiggwilson/songtool/pkg/cmd/internal/config"
)

type ConfigEditCmd struct{}

func (cmd *ConfigEditCmd) Run(cfg *config.Config) error {
	args := append([]string{cfg.Path}, cfg.Edit.Args...)
	exe := exec.Command(cfg.Edit.Command, args...)
	return exe.Run()
}
