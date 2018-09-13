package encrypt

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"

	"github.com/urfave/cli"
)

var encryptRegistryCommand = cli.Command{
	Name:      "registry",
	Usage:     "encrypt registry credentials",
	ArgsUsage: "<repo/name> <string>",
	Action:    encryptRegistry,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "username",
			Usage: "registry username",
		},
		cli.StringFlag{
			Name:  "password",
			Usage: "registry password",
		},
		cli.StringFlag{
			Name:  "server",
			Usage: "registry server",
		},
	},
}

func encryptRegistry(c *cli.Context) error {
	repo := c.Args().First()
	owner, name, err := internal.ParseRepo(repo)
	if err != nil {
		return err
	}

	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	plaintext := c.Args().Get(1)
	if strings.HasPrefix(plaintext, "@") {
		data, err := ioutil.ReadFile(plaintext)
		if err != nil {
			return err
		}
		plaintext = string(data)
	}

	secret := &drone.Secret{
		Data: plaintext,
		Pull: c.Bool("allow-pull-request"),
	}
	encrypted, err := client.EncryptSecret(owner, name, secret)
	if err != nil {
		return err
	}
	fmt.Println(encrypted)
	return nil
}
