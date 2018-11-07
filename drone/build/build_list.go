package build

import (
	"os"
	"text/template"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/urfave/cli"
)

var buildListCmd = cli.Command{
	Name:      "ls",
	Usage:     "show build history",
	ArgsUsage: "<repo/name>",
	Action:    buildList,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplBuildList,
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
		cli.IntFlag{
			Name:  "limit",
			Usage: "limit the list size",
			Value: 25,
		},
		cli.IntFlag{
			Name:  "page",
			Usage: "size of the page",
			Value: 1,
		},
	},
}

func buildList(c *cli.Context) error {
	repo := c.Args().First()
	owner, name, err := internal.ParseRepo(repo)
	if err != nil {
		return err
	}

	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	page := c.Int("page")

	builds, err := client.BuildList(owner, name, page)
	if err != nil {
		return err
	}

	tmpl, err := template.New("_").Parse(c.String("format") + "\n")
	if err != nil {
		return err
	}

	branch := c.String("branch")
	event := c.String("event")
	status := c.String("status")
	limit := c.Int("limit")

	var count int
	for _, build := range builds {
		if count >= limit {
			break
		}
		if branch != "" && build.Target != branch {
			continue
		}
		if event != "" && build.Event != event {
			continue
		}
		if status != "" && build.Status != status {
			continue
		}
		tmpl.Execute(os.Stdout, build)
		count++
	}
	return nil
}

// template for build list information
var tmplBuildList = "\x1b[33mBuild #{{ .Number }} \x1b[0m" + `
Status: {{ .Status }}
Event: {{ .Event }}
Commit: {{ .After }}
Branch: {{ .Target }}
Ref: {{ .Ref }}
Author: {{ .Author }} {{ if .AuthorEmail }}<{{.AuthorEmail}}>{{ end }}
Message: {{ .Message }}
`
