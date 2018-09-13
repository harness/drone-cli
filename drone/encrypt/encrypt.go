package encrypt

import "github.com/urfave/cli"

// Command exports the build command set.
var Command = cli.Command{
	Name:  "encrypt",
	Usage: "encrypt resources",
	Subcommands: []cli.Command{
		encryptSecretCommand,
		encryptRegistryCommand,
	},
}
