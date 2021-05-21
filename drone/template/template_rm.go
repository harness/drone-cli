package template

import (
	"github.com/drone/drone-cli/drone/internal"
	"github.com/urfave/cli"
)

var templateDeleteCmd = cli.Command{
	Name:      "rm",
	Usage:     "remove a template",
	ArgsUsage: "[name]",
	Action:    templateDelete,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "template name",
		},
	},
}

func templateDelete(c *cli.Context) error {
	var (
		name = c.String("name")
	)
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	return client.TemplateDelete(name)
}
