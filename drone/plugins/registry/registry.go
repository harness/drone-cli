package registry

import "github.com/urfave/cli"

// Command exports the registry command set.
var Command = cli.Command{
	Name:  "registry",
	Usage: "registry plugin helpers",
	Subcommands: []cli.Command{
		registryListCmd,
	},
}
