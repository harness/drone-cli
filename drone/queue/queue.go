package queue

import "github.com/urfave/cli"

// Command exports the queue command set.
var Command = cli.Command{
	Name:  "queue",
	Usage: "queue operations",
	Subcommands: []cli.Command{
		queueListCmd,
		queuePauseCmd,
		queueResumeCmd,
	},
}
