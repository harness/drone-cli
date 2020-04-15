package build

import (
	"fmt"
	"strconv"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/urfave/cli"
)

var buildStopCmd = cli.Command{
	Name:      "stop",
	Usage:     "stop a build",
	ArgsUsage: "<repo/name> [build]",
	Action:    buildStop,
}

func buildStop(c *cli.Context) (err error) {
	repo := c.Args().First()
	owner, name, err := internal.ParseRepo(repo)
	if err != nil {
		return err
	}
	number, err := strconv.Atoi(c.Args().Get(1))
	if err != nil {
		return err
	}

	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	err = client.BuildCancel(owner, name, number)
	if err != nil {
		return err
	}

	fmt.Printf("Stopping build %s/%s#%d\n", owner, name, number)
	return nil
}
