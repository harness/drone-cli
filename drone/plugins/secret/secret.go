package secret

import "github.com/urfave/cli"

// Command exports the registry command set.
var Command = cli.Command{
	Name:  "secret",
	Usage: "secret plugin helpers",
	Subcommands: []cli.Command{
		secretFindCmd,
	},
}
