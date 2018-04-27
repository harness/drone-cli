package autoscale

import (
	"github.com/urfave/cli"

	"github.com/drone/drone-cli/drone/internal"
)

var autoscaleResumeCmd = cli.Command{
	Name:   "resume",
	Usage:  "resume the autoscaler",
	Action: autoscaleResume,
}

func autoscaleResume(c *cli.Context) error {
	client, err := internal.NewAutoscaleClient(c)
	if err != nil {
		return err
	}
	return client.AutoscaleResume()
}
