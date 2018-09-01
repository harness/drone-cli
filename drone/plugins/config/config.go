package config

import "github.com/urfave/cli"

// Command exports the registry command set.
var Command = cli.Command{
	Name:  "config",
	Usage: "config plugin helpers",
	Subcommands: []cli.Command{
		configFindCmd,
	},
}
