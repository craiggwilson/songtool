package cmd

type ChordsCmd struct {
	Parse ChordsParseCmd `cmd:"" help:"Parse a chord for validity and proper naming."`
}
