package repo

import (
	"os"
	"text/template"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/funcmap"
	"github.com/urfave/cli"
)

var repoListCmd = cli.Command{
	Name:      "ls",
	Usage:     "list all repos",
	ArgsUsage: " ",
	Action:    repoList,
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
		cli.BoolFlag{
			Name:  "active",
			Usage: "filter active repositories only",
		},
	},
}

func repoList(c *cli.Context) error {
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	repos, err := client.RepoList()
	if err != nil || len(repos) == 0 {
		return err
	}

	tmpl, err := template.New("_").Funcs(funcmap.Funcs).Parse(c.String("format") + "\n")
	if err != nil {
		return err
	}

	org, active := c.String("org"), c.Bool("active")
	for _, repo := range repos {
		if org != "" && org != repo.Namespace {
			continue
		}
		if !repo.Active && active {
			continue
		}
		tmpl.Execute(os.Stdout, repo)
	}
	return nil
}

// template for repository list items
var tmplRepoList = `{{ .Slug }}`
