package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"

	"github.com/codegangsta/cli"
	"github.com/drone/drone-go/drone"
)

var SignCmd = cli.Command{
	Name:  "sign",
	Usage: "creates a secure yaml file",
	Action: func(c *cli.Context) {
		handle(c, signCmd)
	},
	Flags: []cli.Flag{},
}

func signCmd(c *cli.Context, client drone.Client) error {
	repo := c.Args().First()
	owner, name, err := parseRepo(repo)
	if err != nil {
		return err
	}

	in, err := readInput(".drone.yml")
	if err != nil {
		return err
	}

	checksum := shasum(in)
	sig, err := client.Sign(owner, name, []byte(checksum))
	if err != nil {
		return err
	}

	return ioutil.WriteFile(".drone.yml.sig", sig, 0664)
}

func shasum(in []byte) string {
	h := sha256.New()
	h.Write(in)
	return fmt.Sprintf("%x", h.Sum(nil))
}
