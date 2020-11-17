package cron

import (
	"errors"

	"github.com/drone/drone-cli/drone/internal"

	"github.com/urfave/cli"
)

var cronExecCmd = cli.Command{
	Name:      "exec",
	Usage:     "exec cron jobs",
	ArgsUsage: "[repo/name] [cronjob]",
	Action:    cronExec,
}

func cronExec(c *cli.Context) error {
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
	if cron == "" {
		return errors.New("missing cronjob name")
	}

	err = client.CronExec(owner, name, cron)
	return err
}
