package encrypt

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/drone/drone-cli/drone/internal"

	"github.com/urfave/cli"
)

// Command exports the deploy command.
var Command = cli.Command{
	Name:      "encrypt",
	Usage:     "encrypt a string",
	ArgsUsage: "<repo/name> <string>",
	Action:    encrypt,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "aesgcm",
			Usage: "aesgcm encryption",
		},
		cli.BoolFlag{
			Name:  "secretbox",
			Usage: "secretbox encryption",
		},
	},
}

func encrypt(c *cli.Context) error {
	repo := c.Args().First()
	owner, name, err := internal.ParseRepo(repo)
	if err != nil {
		return err
	}

	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	var algorithm string
	switch {
	case c.Bool("asesgcm"):
		algorithm = "asesgcm"
	default:
		algorithm = "secretbox"
	}

	plaintext := c.Args().Get(1)
	if strings.HasPrefix(plaintext, "@") {
		data, err := ioutil.ReadFile(plaintext)
		if err != nil {
			return err
		}
		plaintext = string(data)
	}

	ciphertext, err := client.Encrypt(owner, name, plaintext, algorithm)
	if err != nil {
		return err
	}
	fmt.Println(ciphertext)
	return nil
}
