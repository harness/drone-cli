package main

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/drone/drone-go/drone"
)

var SecretCmd = cli.Command{
	Name:  "secret",
	Usage: "manage secrets",
	Subcommands: []cli.Command{
		// Secret Add
		{
			Name:  "add",
			Usage: "add a secret",
			Action: func(c *cli.Context) {
				handle(c, SecretAddCmd)
			},
		},
		// Secret Delete
		{
			Name:  "rm",
			Usage: "remove a secret",
			Action: func(c *cli.Context) {
				handle(c, SecretDelCmd)
			},
		},
	},
}

func SecretAddCmd(c *cli.Context, client drone.Client) error {
	repo := c.Args().First()
	owner, name, err := parseRepo(repo)
	if err != nil {
		return err
	}

	in := c.Args().Get(1)
	kv := strings.SplitN(in, "=", 2)
	if len(kv) != 2 {
		return fmt.Errorf("Please define the secret in KEY=VALUE format")
	}

	secret := &drone.Secret{}
	secret.Name = kv[0]
	secret.Value = kv[1]
	secret.Image = []string{}
	secret.Event = []string{}
	return client.SecretPost(owner, name, secret)
}

func SecretDelCmd(c *cli.Context, client drone.Client) error {
	repo := c.Args().First()
	owner, name, err := parseRepo(repo)
	if err != nil {
		return err
	}

	secret := c.Args().Get(1)
	return client.SecretDel(owner, name, secret)
}
