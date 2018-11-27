package registry

import (
	"context"
	"os"
	"text/template"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/registry"

	"github.com/urfave/cli"
)

var registryListCmd = cli.Command{
	Name:   "list",
	Usage:  "list the registry credentials",
	Action: registryList,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "ref",
			Usage: "git reference",
			Value: "refs/heads/master",
		},
		cli.StringFlag{
			Name:  "source",
			Usage: "source branch",
		},
		cli.StringFlag{
			Name:  "target",
			Usage: "target branch",
		},
		cli.StringFlag{
			Name:  "before",
			Usage: "commit sha before the change",
		},
		cli.StringFlag{
			Name:  "after",
			Usage: "commit sha after the change",
		},
		cli.StringFlag{
			Name:  "event",
			Usage: "build event",
		},
		cli.StringFlag{
			Name:  "repo",
			Usage: "repository name",
		},

		// TODO(bradrydzewski) these parameters should
		// be defined globally for all plugin commands.

		cli.StringFlag{
			Name:   "endpoint",
			Usage:  "plugin endpoint",
			EnvVar: "DRONE_SECRET_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "secret",
			Usage:  "plugin secret",
			EnvVar: "DRONE_SECRET_SECRET",
		},
		cli.StringFlag{
			Name:   "ssl-skip-verify",
			Usage:  "plugin ssl verification disabled",
			EnvVar: "DRONE_SECRET_VERIFY",
		},
		cli.StringFlag{
			Name:   "format",
			Value:  tmplList,
			Hidden: true,
		},
	},
}

func registryList(c *cli.Context) error {
	slug := c.String("repo")
	owner, name, _ := internal.ParseRepo(slug)

	repo := drone.Repo{
		Namespace: owner,
		Name:      name,
		Slug:      slug,
	}

	build := drone.Build{
		Ref:    c.String("ref"),
		Before: c.String("before"),
		After:  c.String("after"),
		Source: c.String("source"),
		Target: c.String("target"),
		Event:  c.String("event"),
	}

	req := &registry.Request{
		Repo:  repo,
		Build: build,
	}

	client := registry.Client(
		c.String("endpoint"),
		c.String("secret"),
		c.Bool("ssl-skip-verify"),
	)
	list, err := client.List(context.Background(), req)
	if err != nil {
		return err
	}

	format := c.String("format") + "\n"
	tmpl, err := template.New("_").Parse(format)
	if err != nil {
		return err
	}
	for _, item := range list {
		tmpl.Execute(os.Stdout, item)
	}
	return nil
}

var tmplList = "\x1b[33m{{ .Address }} \x1b[0m" + `
Username:  {{ .Username }}
Password: {{ .Password }}
`
