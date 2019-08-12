package user

import (
	"fmt"

	"github.com/drone/drone-go/drone"
	"github.com/urfave/cli"

	"github.com/drone/drone-cli/drone/internal"
)

var userAddCmd = cli.Command{
	Name:      "add",
	Usage:     "adds a user",
	ArgsUsage: "<username>",
	Action:    userAdd,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "admin",
			Usage: "admin privileged",
		},
		cli.BoolFlag{
			Name:  "machine",
			Usage: "machine account",
		},
		cli.StringFlag{
			Name:  "token",
			Usage: "api token",
		},
	},
}

func userAdd(c *cli.Context) error {
	login := c.Args().First()

	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	if login == "" {
		err := fmt.Errorf("Error: blank username not allowed.")
		return err
	}

	in := &drone.User{
		Login:   login,
		Admin:   c.Bool("admin"),
		Machine: c.Bool("machine"),
		Token:   c.String("token"),
	}

	user, err := client.UserCreate(in)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully added user %s\n", user.Login)
	if user.Token != "" {
		fmt.Printf("Generated account token %s\n", user.Token)
	}
	return nil
}
