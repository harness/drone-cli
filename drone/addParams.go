package main

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/drone/drone-go/drone"
	yaml "gopkg.in/yaml.v1"
)

// NewAddParamCommand returns the CLI command for "add-params".
func NewAddParamCommand() cli.Command {
	return cli.Command{
		Name:  "add-params",
		Usage: `adds params for the repo, provide a comma seperated list of key/value pairs ("key1: 'val1',key2: 'val2')`,
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			handle(c, addParamCommandFunc)
		},
	}
}

// addParamCommandFunc executes the "add-params" command.
func addParamCommandFunc(c *cli.Context, client *drone.Client) error {
	var host, owner, name string
	var args = c.Args()
	host, owner, name = parseRepo(args)

	repo, err := client.Repos.Get(host, owner, name)
	if err != nil {
		return fmt.Errorf("Failed to get existing repo params: %s", err)
	}

	// unmarshal the existing parameters into a map
	var currentParams = make(map[string]interface{})
	err = yaml.Unmarshal([]byte(repo.Params), currentParams)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal existing params: %s", err)
	}

	// new parameters will be in the last argument
	var newParamsStr = strings.Replace(args[len(args)-1], ",", "\n", -1)

	var newParams = make(map[string]interface{})
	err = yaml.Unmarshal([]byte(newParamsStr), newParams)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal new params: %s", err)
	}

	// add the new params
	for k, v := range newParams {
		currentParams[k] = v
	}

	updated, err := yaml.Marshal(currentParams)
	if err != nil {
		return fmt.Errorf("Failed to marshal new updated params: %s", err)
	}

	return client.Repos.SetParams(host, owner, name, string(updated))
}
