package template

import (
	"github.com/drone/drone-cli/drone/internal"
	"github.com/urfave/cli"
)

var templateDeleteCmd = cli.Command{
	Name:      "rm",
	Usage:     "remove a template",
	ArgsUsage: "[namespace] [name]",
	Action:    templateDelete,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "template name",
		},
		cli.StringFlag{
			Name:  "namespace",
			Usage: "organization name",
		},
	},
}

func templateDelete(c *cli.Context) error {
	var (
		namespace = c.String("namespace")
		name      = c.String("name")
	)
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	return client.TemplateDelete(namespace, name)
}
