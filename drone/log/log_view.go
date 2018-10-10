package log

import (
	"strconv"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/urfave/cli"
)

var logViewCmd = cli.Command{
	Name:      "view",
	Usage:     "display the step logs",
	ArgsUsage: "<repo/name> <build> <stage> <step>",
	Action:    logView,
}

func logView(c *cli.Context) (err error) {
	repo := c.Args().First()
	owner, name, err := internal.ParseRepo(repo)
	if err != nil {
		return err
	}
	number, err := strconv.Atoi(c.Args().Get(1))
	if err != nil {
		return err
	}
	stage, err := strconv.Atoi(c.Args().Get(2))
	if err != nil {
		return err
	}
	step, err := strconv.Atoi(c.Args().Get(3))
	if err != nil {
		return err
	}

	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	lines, err := client.Logs(owner, name, number, stage, step)
	if err != nil {
		return err
	}

	for _, line := range lines {
		print(line.Message)
	}
	return nil
}
