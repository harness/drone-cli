package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/codegangsta/cli"
)

var repoCmd = cli.Command{
	Name:  "repo",
	Usage: "manage repositories",
	Subcommands: []cli.Command{

		// list command
		cli.Command{
			Name:   "ls",
			Usage:  "list all repos",
			Action: repoList,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "format",
					Usage: "format output",
					Value: tmplRepoList,
				},
				cli.StringFlag{
					Name:  "org",
					Usage: "filter by organization",
				},
			},
		},

		// info command
		cli.Command{
			Name:   "info",
			Usage:  "show repository details",
			Action: repoInfo,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "format",
					Usage: "format output",
					Value: tmplRepoInfo,
				},
			},
		},

		// add command
		cli.Command{
			Name:   "add",
			Usage:  "add a repository",
			Action: repoAdd,
		},

		// remove command
		cli.Command{
			Name:   "rm",
			Usage:  "remove a repository",
			Action: repoRemove,
		},
	},
}

// command to fetch and return repository information.
func repoInfo(c *cli.Context) error {
	arg := c.Args().First()
	owner, name, err := parseRepo(arg)
	if err != nil {
		return err
	}

	client, err := newClient(c)
	if err != nil {
		return err
	}

	repo, err := client.Repo(owner, name)
	if err != nil {
		return err
	}

	tmpl, err := template.New("_").Parse(c.String("format"))
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, repo)
}

// command to add a repository.
func repoAdd(c *cli.Context) error {
	repo := c.Args().First()
	owner, name, err := parseRepo(repo)
	if err != nil {
		return err
	}

	client, err := newClient(c)
	if err != nil {
		return err
	}

	if _, err := client.RepoPost(owner, name); err != nil {
		return err
	}
	fmt.Printf("Successfully activated repository %s/%s\n", owner, name)
	return nil
}

// command to remove a repository from drone.
func repoRemove(c *cli.Context) error {
	repo := c.Args().First()
	owner, name, err := parseRepo(repo)
	if err != nil {
		return err
	}

	client, err := newClient(c)
	if err != nil {
		return err
	}

	if err := client.RepoDel(owner, name); err != nil {
		return err
	}
	fmt.Printf("Successfully removed repository %s/%s\n", owner, name)
	return nil
}

// command to list user repositories.
func repoList(c *cli.Context) error {
	client, err := newClient(c)
	if err != nil {
		return err
	}

	repos, err := client.RepoList()
	if err != nil || len(repos) == 0 {
		return err
	}

	tmpl, err := template.New("_").Parse(c.String("format") + "\n")
	if err != nil {
		return err
	}

	org := c.String("org")
	for _, repo := range repos {
		if org != "" && org != repo.Owner {
			continue
		}
		tmpl.Execute(os.Stdout, repo)
	}
	return nil
}

// repository info template.
var tmplRepoInfo = `Owner: {{ .Owner }}
Repo: {{ .Name }}
Type: {{ .Kind }}
Private: {{ .IsPrivate }}
Remote: {{ .Clone }}
`

// repository list template.
var tmplRepoList = `{{ .FullName }}`
