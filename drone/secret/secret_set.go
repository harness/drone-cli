package secret

import (
	"github.com/drone/drone-cli/drone/internal"

	"github.com/urfave/cli"
)

var secretUpdateCmd = cli.Command{
	Name:      "update",
	Usage:     "update a secret",
	ArgsUsage: "[repo/name]",
	Action:    secretUpdate,
	Flags:     flags,
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
	secret, err := makeSecret(c)
	if err != nil {
		return err
	}
	_, err = client.SecretUpdate(owner, name, secret)
	if err == nil {
		return printSecret(secret, tmplSecretList)
	}
	return err
}
