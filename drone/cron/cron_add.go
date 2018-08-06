package cron

import (
	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"

	"github.com/urfave/cli"
)

var cronCreateCmd = cli.Command{
	Name:      "add",
	Usage:     "adds a cronjob",
	ArgsUsage: "[repo/name]",
	Action:    cronCreate,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "branch",
			Usage: "branch name",
			Value: "master",
		},
	},
}

func cronCreate(c *cli.Context) error {
	slug := c.Args().First()
	owner, name, err := internal.ParseRepo(slug)
	if err != nil {
		return err
	}
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	cron := &drone.Cron{
		Name:   c.Args().Get(1),
		Expr:   c.Args().Get(2),
		Branch: c.String("branch"),
	}
	_, err = client.CronCreate(owner, name, cron)
	if err != nil {
		return err
	}
	return nil
}
