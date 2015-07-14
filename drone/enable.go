package main

import (
	"github.com/codegangsta/cli"
	"github.com/drone/drone-go/drone"
)

// NewEnableCommand returns the CLI command for "enable".
func NewEnableCommand() cli.Command {
	return cli.Command{
		Name:  "enable",
		Usage: "enable a repository",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			handle(c, enableCommandFunc)
		},
	}
}

// enableCommandFunc executes the "enable" command.
func enableCommandFunc(c *cli.Context, client *drone.Client) error {
	host, owner, name := parseRepo(c.Args())
	return client.Repos.Enable(host, owner, name)
}
