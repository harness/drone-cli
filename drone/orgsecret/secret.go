package orgsecret

import "github.com/urfave/cli"

// Command exports the secret command.
var Command = cli.Command{
	Name:  "orgsecret",
	Usage: "manage organization secrets",
	Subcommands: []cli.Command{
		secretCreateCmd,
		secretDeleteCmd,
		secretUpdateCmd,
		secretInfoCmd,
		secretListCmd,
	},
}
