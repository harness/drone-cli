package build

import "github.com/urfave/cli"

// Command exports the build command set.
var Command = cli.Command{
	Name:  "build",
	Usage: "manage builds",
	Subcommands: []cli.Command{
		buildListCmd,
		buildLastCmd,
		buildInfoCmd,
		buildStopCmd,
		buildStartCmd,
		buildApproveCmd,
		buildDeclineCmd,
		buildPromoteCmd,
		buildRollbackCmd,
		buildQueueCmd,
	},
}
