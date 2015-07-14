package main

import (
	"github.com/codegangsta/cli"
	"github.com/drone/drone-go/drone"
)

// NewDeleteCommand returns the CLI command for "delete".
func NewDeleteCommand() cli.Command {
	return cli.Command{
		Name:  "delete",
		Usage: "delete a repository",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			handle(c, deleteCommandFunc)
		},
	}
}

// deleteCommandFunc executes the "delete" command.
func deleteCommandFunc(c *cli.Context, client *drone.Client) error {
	var host, owner, name string
	host, owner, name = parseRepo(c.Args())

	return client.Repos.Delete(host, owner, name)
}
