package main

import (
	"fmt"
	"io/ioutil"

	"github.com/codegangsta/cli"
	"github.com/drone/drone-go/drone"
)

// NewSetKeyCommand returns the CLI command for "set-key".
func NewSetKeyCommand() cli.Command {
	return cli.Command{
		Name:  "set-key",
		Usage: "sets the SSH private key used to clone, pass the path of the private key",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			handle(c, setKeyCommandFunc)
		},
	}
}

// setKeyCommandFunc executes the "set-key" command.
func setKeyCommandFunc(c *cli.Context, client *drone.Client) error {
	var keyPath string
	var args = c.Args()
	var host, owner, name = parseRepo(c.Args())

	// path to private key will be last arg
	if len(args) == 0 {
		return fmt.Errorf("The path to the private key is required.")
	}
	keyPath = args[len(args)-1]

	privKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("Could not find private RSA key %s. %s", keyPath, err)
	}

	pathPub := keyPath + ".pub"
	pubKey, err := ioutil.ReadFile(pathPub)
	if err != nil {
		return fmt.Errorf("Could not find public RSA key %s. %s", pathPub, err)
	}

	return client.Repos.SetKey(host, owner, name, string(pubKey), string(privKey))
}
