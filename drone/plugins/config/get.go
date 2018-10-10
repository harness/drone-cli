package config

import (
	"context"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/config"
	"github.com/urfave/cli"
)

var configFindCmd = cli.Command{
	Name:      "get",
	Usage:     "get the pipeline configuration",
	ArgsUsage: "[repo/name]",
	Action:    configFind,
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

		// TODO(bradrydzewski) these parameters should
		// be defined globally for all plugin commands.

		cli.StringFlag{
			Name:   "endpoint",
			Usage:  "plugin endpoint",
			EnvVar: "DRONE_YAML_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "secret",
			Usage:  "plugin secret",
			EnvVar: "DRONE_YAML_SECRET",
		},
		cli.StringFlag{
			Name:   "ssl-skip-verify",
			Usage:  "plugin ssl verification disabled",
			EnvVar: "DRONE_YAML_SKIP_VERIFY",
		},
	},
}

func configFind(c *cli.Context) error {
	slug := c.Args().First()
	owner, name, err := internal.ParseRepo(slug)
	if err != nil {
		return err
	}

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
	}

	req := &config.Request{
		Repo:  repo,
		Build: build,
	}

	client := config.Client(
		c.String("endpoint"),
		c.String("secret"),
		c.Bool("ssl-skip-verify"),
	)
	res, err := client.Find(context.Background(), req)
	if err != nil {
		return err
	}
	println(res.Data)
	return nil
}
