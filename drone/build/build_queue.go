package build

import (
	"os"
	"text/template"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/funcmap"
	"github.com/urfave/cli"
)

var buildQueueCmd = cli.Command{
	Name:      "queue",
	Usage:     "show build queue",
	ArgsUsage: "",
	Action:    buildQueue,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplBuildQueue,
		},
		cli.StringFlag{
			Name:  "repo",
			Usage: "repo filter",
		},
		cli.StringFlag{
			Name:  "branch",
			Usage: "branch filter",
		},
		cli.StringFlag{
			Name:  "event",
			Usage: "event filter",
		},
		cli.StringFlag{
			Name:  "status",
			Usage: "status filter",
		},
	},
}

func buildQueue(c *cli.Context) error {
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	repos, err := client.Incomplete()
	if err != nil {
		return err
	}

	tmpl, err := template.New("_").Funcs(funcmap.Funcs).Parse(c.String("format") + "\n")
	if err != nil {
		return err
	}

	slug := c.String("repo")
	branch := c.String("branch")
	event := c.String("event")
	status := c.String("status")

	for _, repo := range repos {
		if slug != "" && repo.Slug != slug {
			continue
		}
		if branch != "" && repo.Build.Target != branch {
			continue
		}
		if event != "" && repo.Build.Event != event {
			continue
		}
		if status != "" && repo.Build.Status != status {
			continue
		}
		tmpl.Execute(os.Stdout, repo)
	}
	return nil
}

// template for build queue information
var tmplBuildQueue = "\x1b[33m{{ .Slug }}#{{ .Build.Number }} \x1b[0m" + `
Name: {{ .Slug }}
Build: {{ .Build.Number }}
Status: {{ .Build.Status }}
Event: {{ .Build.Event }}
Branch: {{ .Build.Target }}
Ref: {{ .Build.Ref }}
Author: {{ .Build.Author }}{{ if .Build.AuthorEmail }} <{{ .Build.AuthorEmail }}>{{ end }}
Created: {{ .Build.Created | time }}
Started: {{ .Build.Started | time }}
Updated: {{ .Build.Updated | time }}
`
