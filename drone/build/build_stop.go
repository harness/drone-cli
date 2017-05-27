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
	ArgsUsage: "<repo/name> [build] [job]",
	Action:    buildStop,
}

func buildStop(c *cli.Context) (err error) {
	repo := c.Args().First()
	owner, name, err := internal.ParseRepo(repo)
	if err != nil {
		return err
	}
	number, err := parseBuildArg(c.Args().Get(1))
	if err != nil {
		return err
	}
	var job int

	jobIdStr := c.Args().Get(2)
	if len(jobIdStr) == 0 {
		job = 1
	} else {
		job, err = strconv.Atoi(jobIdStr)
		if err != nil {
			return fmt.Errorf("Error: malformed job number specified: %s", err.Error())
		}
	}

	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	err = client.BuildStop(owner, name, number, job)
	if err != nil {
		return err
	}

	fmt.Printf("Stopping build %s/%s#%d.%d\n", owner, name, number, job)
	return nil
}
