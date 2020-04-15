package queue

import (
	"github.com/drone/drone-cli/drone/internal"
	"github.com/urfave/cli"
)

var queueResumeCmd = cli.Command{
	Name:   "resume",
	Usage:  "resume queue operations",
	Action: queueResume,
}

func queueResume(c *cli.Context) (err error) {
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	return client.QueueResume()
}
