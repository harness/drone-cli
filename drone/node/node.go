package node

import "github.com/urfave/cli"

// Command exports the registry command set.
var Command = cli.Command{
	Name:   "node",
	Usage:  "manage nodes",
	Hidden: true,
	Subcommands: []cli.Command{
		nodeListCmd,
		nodeInfoCmd,
		nodeCreateCmd,
		// nodeUpdateCmd,
		// nodeDeleteCmd,
		// nodePauseCmd,
		// nodeUnpauseCmd,
		// nodeLockCmd,
		// nodeUnlockCmd,
		// nodeInitCmd,
		nodeImportCmd,
		// nodeKeygenCmd,
	},
}
