package plugins

import (
	"github.com/drone/drone-cli/drone/plugins/config"

	"github.com/urfave/cli"
)

// Command exports the registry command set.
var Command = cli.Command{
	Name:  "plugins",
	Usage: "plugin helper functions",
	Subcommands: []cli.Command{
		config.Command,
	},
}
