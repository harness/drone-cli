package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/drone/drone-go/drone"
)

// NewGetParamsCommand returns the CLI command for "get-params".
func NewGetParamsCommand() cli.Command {
	return cli.Command{
		Name:  "get-params",
		Usage: "gets params for the repo",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			handle(c, getParamsCommandFunc)
		},
	}
}

// getParamsCommandFunc executes the "get-params" command.
func getParamsCommandFunc(c *cli.Context, client *drone.Client) error {
	host, owner, name := parseRepo(c.Args())

	repo, err := client.Repos.Get(host, owner, name)
	if err != nil {
		return err
	}

	fmt.Printf("%s/%s/%s Parameters:\n%s\n", host, owner, name, repo.Params)
	return nil
}
