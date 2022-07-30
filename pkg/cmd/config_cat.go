package cmd

type ConfigCatCmd struct {
}

func (cmd *ConfigCatCmd) Run(cfg *Config) error {
	return printJSON(cfg.ConfigFile)
}
