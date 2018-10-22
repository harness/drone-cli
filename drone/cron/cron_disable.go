package cron

import (
	"errors"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"

	"github.com/urfave/cli"
)

var cronDisableCmd = cli.Command{
	Name:      "disable",
	Usage:     "disable cron jobs",
	ArgsUsage: "[repo/name] [cronjob]",
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
	cron := c.Args().Get(1)
	if cron == "" {
		return errors.New("missing cronjob name")
	}
	disabled := true
	in := &drone.CronPatch{
		Disabled: &disabled,
	}
	_, err = client.CronUpdate(owner, name, cron, in)
	return err
}
