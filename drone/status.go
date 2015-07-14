package main

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/drone/drone-go/drone"
)

// NewStatusCommand returns the CLI command for "status".
func NewStatusCommand() cli.Command {
	return cli.Command{
		Name:  "status",
		Usage: "display a repository build status",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "b, branch",
				Usage: "branch to display",
			},
		},
		Action: func(c *cli.Context) {
			handle(c, statusCommandFunc)
		},
	}
}

// statusCommandFunc executes the "status" command.
func statusCommandFunc(c *cli.Context, client *drone.Client) error {
	host, owner, repo := parseRepo(c.Args())
	branch := "master"

	if c.IsSet("branch") {
		branch = c.String("branch")
	}

	builds, err := client.Commits.List(host, owner, repo)
	if err != nil {
		return err
	} else if len(builds) == 0 {
		log.Printf("No builds found for %s/%s/%s", host, owner, repo)
		return nil
	}

	// builds go from older to newer, we want to grab the newest build
	for i := len(builds) - 1; i >= 0; i-- {
		if builds[i].Branch == branch {
			fmt.Printf("%s\t%s\t%s\t%s\t%v\n", builds[i].Status, builds[i].ShaShort(), builds[i].Timestamp, builds[i].Author, builds[i].Message)
			return nil
		}
	}

	return fmt.Errorf("Could not find builds on branch: %s", branch)
}
