package cron

import (
	"errors"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"

	"github.com/urfave/cli"
)

var cronEnableCmd = cli.Command{
	Name:      "enable",
	Usage:     "enable cron jobs",
	ArgsUsage: "[repo/name] [cronjob]",
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
	cron := c.Args().Get(1)
	if cron == "" {
		return errors.New("missing cronjob name")
	}
	disabled := false
	in := &drone.CronPatch{
		Disabled: &disabled,
	}
	_, err = client.CronUpdate(owner, name, cron, in)
	return err
}
