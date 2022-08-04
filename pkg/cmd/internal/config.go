package internal

type ConfigCmd struct {
	Cat  ConfigCatCmd  `cmd:"" help:"Prints the config."`
	Edit ConfigEditCmd `cmd:"" help:"Launches an editor for the config file."`
	Path ConfigPathCmd `cmd:"" help:"Prints the location of the config file."`
}
