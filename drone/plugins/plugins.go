package plugins

import (
	"github.com/drone/drone-cli/drone/plugins/admit"
	"github.com/drone/drone-cli/drone/plugins/config"
	"github.com/drone/drone-cli/drone/plugins/convert"
	"github.com/drone/drone-cli/drone/plugins/registry"
	"github.com/drone/drone-cli/drone/plugins/secret"

	"github.com/urfave/cli"
)

// Command exports the registry command set.
var Command = cli.Command{
	Name:  "plugins",
	Usage: "plugin helper functions",
	Subcommands: []cli.Command{
		admit.Command,
		config.Command,
		convert.Command,
		secret.Command,
		registry.Command,
	},
}
