package convert

import (
	"context"
	"io/ioutil"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
	"github.com/urfave/cli"
)

// Command exports the registry command set.
var Command = cli.Command{
	Name:      "convert",
	Usage:     "convert the pipeline configuration",
	ArgsUsage: "[repo/name]",
	Action:    convert,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Usage: "path to the configuration file",
		},
		cli.StringFlag{
			Name:  "ref",
			Usage: "git reference",
			Value: "refs/heads/master",
		},
		cli.StringFlag{
			Name:  "source",
			Usage: "source branch",
			Value: "master",
		},
		cli.StringFlag{
			Name:  "target",
			Usage: "target branch",
			Value: "master",
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
			Name:  "repository",
			Usage: "repository name",
		},

		// TODO(bradrydzewski) these parameters should
		// be defined globally for all plugin commands.

		cli.StringFlag{
			Name:   "endpoint",
			Usage:  "plugin endpoint",
			EnvVar: "DRONE_CONVERT_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "secret",
			Usage:  "plugin secret",
			EnvVar: "DRONE_CONVERT_SECRET",
		},
		cli.StringFlag{
			Name:   "ssl-skip-verify",
			Usage:  "plugin ssl verification disabled",
			EnvVar: "DRONE_CONVERT_SKIP_VERIFY",
		},
	},
}

func convert(c *cli.Context) error {
	slug := c.String("repository")
	owner, name, _ := internal.ParseRepo(slug)

	path := c.String("path")
	if path == "" {
		path = c.Args().First()
	}

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	req := &converter.Request{
		Repo: drone.Repo{
			Namespace: owner,
			Name:      name,
			Slug:      slug,
			Config:    path,
		},
		Build: drone.Build{
			Ref:    c.String("ref"),
			Before: c.String("before"),
			After:  c.String("after"),
			Source: c.String("source"),
			Target: c.String("target"),
		},
		Config: drone.Config{
			Data: string(raw),
		},
	}

	client := converter.Client(
		c.String("endpoint"),
		c.String("secret"),
		c.Bool("ssl-skip-verify"),
	)
	res, err := client.Convert(context.Background(), req)
	if err != nil {
		return err
	}
	switch {
	case res == nil:
		println(string(raw))
	case res != nil:
		println(res.Data)
	}
	return nil
}
