package template

import "github.com/urfave/cli"

var Command = cli.Command{
	Name:  "template",
	Usage: "manage templates",
	Subcommands: []cli.Command{
		templateCreateCmd,
		templateInfoCmd,
		templateListCmd,
		templateCreateCmd,
		templateUpdateCmd,
		templateDeleteCmd,
	},
}
