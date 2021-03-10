package cron

import "github.com/urfave/cli"

// Command exports the registry command set.
var Command = cli.Command{
	Name:  "cron",
	Usage: "manage cron jobs",
	Subcommands: []cli.Command{
		cronListCmd,
		cronInfoCmd,
		cronCreateCmd,
		cronDeleteCmd,
		cronDisableCmd,
		cronEnableCmd,
	},
}
