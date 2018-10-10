package sign

import (
	"fmt"
	"io/ioutil"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-yaml/yaml/signer"
	"github.com/urfave/cli"
)

// Command exports the sign command.
var Command = cli.Command{
	Name:      "sign",
	Usage:     "sign the yaml file",
	ArgsUsage: "<source>",
	Action:    format,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "save",
			Usage: "save result to source",
		},
	},
}

func format(c *cli.Context) error {
	repo := c.Args().First()
	owner, name, err := internal.ParseRepo(repo)
	if err != nil {
		return err
	}

	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	path := c.Args().Get(1)
	if path == "" {
		path = ".drone.yml"
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	hmac, err := client.Sign(owner, name, string(data))
	if err != nil {
		return err
	}

	if c.Bool("save") == false {
		fmt.Println(hmac)
		return nil
	}

	data, err = signer.WriteTo(data, hmac)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}
