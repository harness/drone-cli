package cron

import (
	"github.com/drone/drone-cli/drone/internal"

	"github.com/urfave/cli"
)

var cronDeleteCmd = cli.Command{
	Name:      "rm",
	Usage:     "display cron rm",
	ArgsUsage: "[repo/name]",
	Action:    cronDelete,
}

func cronDelete(c *cli.Context) error {
	slug := c.Args().First()
	owner, name, err := internal.ParseRepo(slug)
	if err != nil {
		return err
	}
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	cronjob := c.Args().Get(1)
	return client.CronDelete(owner, name, cronjob)
}
