package sign

import (
	"fmt"
	"github.com/drone/drone-cli/drone/internal"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
)

var Command = cli.Command{
	Name:   "sign",
	Usage:  "sign the yaml file",
	Action: sign,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "in",
			Usage: "input file",
			Value: ".drone.yml",
		},
		cli.StringFlag{
			Name:  "out",
			Usage: "output file signature",
			Value: ".drone.yml.sig",
		},
	},
}

func readInput(in string) ([]byte, error) {
	if in == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(in)
}

func sign(c *cli.Context) error {
	repo := c.Args().First()
	owner, name, err := internal.ParseRepo(repo)
	if err != nil {
		return err
	}

	in, err := readInput(c.String("in"))
	if err != nil {
		return fmt.Errorf("readInput: %v", err)
	}

	client, err := internal.NewClient(c)
	if err != nil {
		return fmt.Errorf("NewClient: %v", err)
	}

	sig, err := client.Sign(owner, name, in)
	if err != nil {
		return fmt.Errorf("Sign: %v", err)
	}

	return ioutil.WriteFile(c.String("out"), sig, 0664)
}
