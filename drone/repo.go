package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/drone/drone-go/drone"
)

var RepoCmd = cli.Command{
	Name:  "repo",
	Usage: "manage repos",
	Subcommands: []cli.Command{
		// Repo List
		{
			Name:  "ls",
			Usage: "lists repositories",
			Action: func(c *cli.Context) {
				handle(c, RepoListCmd)
			},
		},
		// Repo Info
		{
			Name:  "info",
			Usage: "show repository details",
			Action: func(c *cli.Context) {
				handle(c, RepoInfoCmd)
			},
		},
		// Repo Add
		{
			Name:  "add",
			Usage: "add a repository",
			Action: func(c *cli.Context) {
				handle(c, RepoAddCmd)
			},
		},
		// Repo Delete
		{
			Name:  "rm",
			Usage: "remove a repository",
			Action: func(c *cli.Context) {
				handle(c, RepoDelCmd)
			},
		},
	},
}

func RepoAddCmd(c *cli.Context, client drone.Client) error {
	owner, name, err := parseRepo(c.Args().Get(0))
	if err != nil {
		return err
	}

	repo, err := client.RepoPost(owner, name)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully added %s\n", repo.FullName)
	return nil
}

func RepoDelCmd(c *cli.Context, client drone.Client) error {
	owner, name, err := parseRepo(c.Args().Get(0))
	if err != nil {
		return err
	}

	err = client.RepoDel(owner, name)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully removed %s/%s\n", owner, name)
	return nil
}

func RepoListCmd(c *cli.Context, client drone.Client) error {
	repos, err := client.RepoList()
	if err != nil {
		return err
	}

	for _, repo := range repos {
		fmt.Println(repo.FullName)
	}
	return nil
}

func RepoInfoCmd(c *cli.Context, client drone.Client) error {
	owner, name, err := parseRepo(c.Args().Get(0))
	if err != nil {
		return err
	}
	repo, err := client.Repo(owner, name)
	if err != nil {
		return err
	}

	fmt.Println(repo.FullName)
	return nil
}
