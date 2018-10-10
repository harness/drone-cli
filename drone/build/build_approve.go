package build

import (
	"fmt"
	"strconv"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/urfave/cli"
)

var buildApproveCmd = cli.Command{
	Name:      "approve",
	Usage:     "approve a build stage",
	ArgsUsage: "<repo/name> <build> <stage>",
	Action:    buildApprove,
}

func buildApprove(c *cli.Context) (err error) {
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

	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	err = client.Approve(owner, name, number, stage)
	if err != nil {
		return err
	}

	fmt.Printf("Approving build %s/%s#%d\n", owner, name, number)
	return nil
}
