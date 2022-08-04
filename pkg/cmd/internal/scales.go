package internal

type ScalesCmd struct {
	Cat ScalesCatCmd `cmd:"" help:"Prints the scale."`
	Ls  ScalesLsCmd  `cmd:"" help:"Lists the available scales."`
}
