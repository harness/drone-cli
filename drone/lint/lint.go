package lint

import (
	"github.com/drone/drone-yaml/yaml"
	"github.com/drone/drone-yaml/yaml/linter"
	"github.com/urfave/cli"
)

// Command exports the linter command.
var Command = cli.Command{
	Name:      "lint",
	Usage:     "lint the yaml file",
	ArgsUsage: "<source>",
	Action:    lint,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "trusted",
			Usage: "is the yaml trustable",
		},
	},
}

func lint(c *cli.Context) error {
	path := c.Args().First()
	if path == "" {
		path = ".drone.yml"
	}

	manifest, err := yaml.ParseFile(path)
	if err != nil {
		return err
	}

	for _, resource := range manifest.Resources {
		if err := linter.Lint(resource, c.Bool("trusted")); err != nil {
			return err
		}
	}

	return nil
}
