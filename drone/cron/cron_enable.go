package cron

import (
	"github.com/drone/drone-cli/drone/internal"

	"github.com/urfave/cli"
)

var cronEnableCmd = cli.Command{
	Name:      "enable",
	Usage:     "enable cron jobs",
	ArgsUsage: "[repo/name]",
	Action:    cronEnable,
}

func cronEnable(c *cli.Context) error {
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
