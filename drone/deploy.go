package main

import (
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/drone/drone-go/drone"
)

var DeployCmd = cli.Command{
	Name:  "deploy",
	Usage: "deploy code",
	Action: func(c *cli.Context) {
		handle(c, deployCmd)
	},
}

func deployCmd(c *cli.Context, client drone.Client) error {
	var (
		nameParam = c.Args().Get(0)
		numParam  = c.Args().Get(1)

		err   error
		owner string
		name  string
		num   int
	)

	num, err = strconv.Atoi(numParam)
	if err != nil {
		return fmt.Errorf("Invalid or missing build number")
	}
	owner, name, err = parseRepo(nameParam)
	if err != nil {
		return err
	}

	build, err := client.Build(owner, name, num)
	if err != nil {
		return err
	}
	if build.Event == drone.EventPull {
		return fmt.Errorf("Cannot trigger a pull request deployment")
	}

	// client.BuildFork(owner, name, num)

	return nil
}
