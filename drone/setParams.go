package main

import (
	"fmt"
	"io/ioutil"

	"github.com/codegangsta/cli"
	"github.com/drone/drone-go/drone"
)

// NewSetParamsCommand returns the CLI command for "set-params".
func NewSetParamsCommand() cli.Command {
	return cli.Command{
		Name:  "set-params",
		Usage: "sets all params for the repo",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			handle(c, setParamsCommandFunc)
		},
	}
}

// setParamsCommandFunc executes the "set-params" command.
func setParamsCommandFunc(c *cli.Context, client *drone.Client) error {
	var host, owner, name, path string
	var args = c.Args()

	host, owner, name = parseRepo(c.Args())

	if len(args) == 2 {
		path = args[1]
	}

	params, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Could not find params file %s. %s", path, err)
	}

	return client.Repos.SetParams(host, owner, name, string(params))
}
