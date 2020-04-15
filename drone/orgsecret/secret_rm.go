package orgsecret

import (
	"github.com/urfave/cli"

	"github.com/drone/drone-cli/drone/internal"
)

var secretDeleteCmd = cli.Command{
	Name:      "rm",
	Usage:     "remove a secret",
	ArgsUsage: "[organization] [name]",
	Action:    secretDelete,
	Flags:     []cli.Flag{},
}

func secretDelete(c *cli.Context) error {
	var (
		namespace = c.Args().First()
		name      = c.Args().Get(1)
	)
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	return client.OrgSecretDelete(namespace, name)
}
