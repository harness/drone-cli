package cron

import (
	"github.com/drone/drone-cli/drone/internal"

	"github.com/urfave/cli"
)

var cronDisableCmd = cli.Command{
	Name:      "disable",
	Usage:     "disable cron jobs",
	ArgsUsage: "[repo/name]",
	Action:    cronDisable,
}

func cronDisable(c *cli.Context) error {
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
	return client.CronDisable(owner, name, cronjob)
}
