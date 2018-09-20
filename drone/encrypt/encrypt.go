package encrypt

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"

	"github.com/urfave/cli"
)

// Command is an encryption cli.Command
var Command = cli.Command{
	Name:      "encrypt",
	Usage:     "encrypt a secret",
	ArgsUsage: "<repo/name> <string>",
	Action:    encryptSecret,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "allow-pull-request",
			Usage: "permit access to pull requests",
		},
	},
}

func encryptSecret(c *cli.Context) error {
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
	encrypted, err := client.Encrypt(owner, name, secret)
	if err != nil {
		return err
	}
	fmt.Println(encrypted)
	return nil
}
