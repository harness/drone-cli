package secret

import (
	"io/ioutil"
	"strings"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"

	"github.com/urfave/cli"
)

var secretUpdateCmd = cli.Command{
	Name:      "update",
	Usage:     "update a secret",
	ArgsUsage: "[repo/name]",
	Action:    secretUpdate,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "repository",
			Usage: "repository name (e.g. octocat/hello-world)",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "secret name",
		},
		cli.StringFlag{
			Name:  "data",
			Usage: "secret value",
		},
		cli.BoolFlag{
			Name:  "allow-pull-request",
			Usage: "permit read access to pull requests",
		},
		cli.BoolFlag{
			Name:  "allow-push-on-pull-request",
			Usage: "permit write access to pull requests (e.g. allow docker push)",
		},
	},
}

func secretUpdate(c *cli.Context) error {
	reponame := c.String("repository")
	if reponame == "" {
		reponame = c.Args().First()
	}
	owner, name, err := internal.ParseRepo(reponame)
	if err != nil {
		return err
	}
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	secret := &drone.Secret{
		Name:            c.String("name"),
		Data:            c.String("data"),
		PullRequest:     c.Bool("allow-pull-request"),
		PullRequestPush: c.Bool("allow-push-on-pull-request"),
	}
	if strings.HasPrefix(secret.Data, "@") {
		path := strings.TrimPrefix(secret.Data, "@")
		out, ferr := ioutil.ReadFile(path)
		if ferr != nil {
			return ferr
		}
		secret.Data = string(out)
	}
	_, err = client.SecretUpdate(owner, name, secret)
	return err
}
