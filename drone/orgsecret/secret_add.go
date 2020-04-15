package orgsecret

import (
	"io/ioutil"
	"strings"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"

	"github.com/urfave/cli"
)

var secretCreateCmd = cli.Command{
	Name:      "add",
	Usage:     "adds a secret",
	ArgsUsage: "[organization] [name] [data]",
	Action:    secretCreate,
	Flags: []cli.Flag{
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

func secretCreate(c *cli.Context) error {
	var (
		namespace = c.Args().First()
		name      = c.Args().Get(1)
		data      = c.Args().Get(2)
	)
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	secret := &drone.Secret{
		Name:            name,
		Data:            data,
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
	_, err = client.OrgSecretCreate(namespace, secret)
	return err
}
