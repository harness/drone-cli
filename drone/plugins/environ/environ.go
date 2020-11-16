package environ

import (
	"context"
	"fmt"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/environ"
	"github.com/urfave/cli"
)

// Command exports the admission command set.
var Command = cli.Command{
	Name:   "env",
	Usage:  "test env extensions",
	Action: environAction,
	Flags: []cli.Flag{

		//
		// build and repository details
		//

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

		cli.StringFlag{
			Name:   "endpoint",
			Usage:  "plugin endpoint",
			EnvVar: "DRONE_ENVIRON_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "secret",
			Usage:  "plugin secret",
			EnvVar: "DRONE_ENVIRON_SECRET",
		},
		cli.StringFlag{
			Name:   "ssl-skip-verify",
			Usage:  "plugin ssl verification disabled",
			EnvVar: "DRONE_ENVIRON_SKIP_VERIFY",
		},
	},
}

func environAction(c *cli.Context) error {
	slug := c.String("repository")
	owner, name, _ := internal.ParseRepo(slug)

	req := &environ.Request{
		Repo: drone.Repo{
			Namespace: owner,
			Name:      name,
			Slug:      slug,
		},
		Build: drone.Build{
			Ref:    c.String("ref"),
			Before: c.String("before"),
			After:  c.String("after"),
			Source: c.String("source"),
			Target: c.String("target"),
		},
	}

	client := environ.Client(
		c.String("endpoint"),
		c.String("secret"),
		c.Bool("ssl-skip-verify"),
	)
	list, err := client.List(context.Background(), req)
	if err != nil {
		return err
	}

	for k, v := range list {
		fmt.Printf("%s=%s\n", k, v)
	}
	return nil
}
