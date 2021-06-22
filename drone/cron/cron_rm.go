package cron

import (
	"github.com/drone/drone-cli/drone/internal"
	"github.com/urfave/cli"
)

var cronDeleteCmd = cli.Command{
	Name:      "rm",
	Usage:     "deletes a cronjob",
	ArgsUsage: "[repo/name] [cronjob]",
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
	cron := c.Args().Get(1)
	client.CronDelete(owner, name, cron)
	return client.CronDelete(owner, name, cron)
}
