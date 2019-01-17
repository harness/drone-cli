package user

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"
)

var userUpdateCmd = cli.Command{
	Name:      "update",
	Usage:     "update a user",
	ArgsUsage: "<username>",
	Action:    userUpdate,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "active",
			Usage: "user is active",
		},
		cli.BoolFlag{
			Name:  "admin",
			Usage: "user is an admin",
		},
	},
}

func userUpdate(c *cli.Context) error {
	login := c.Args().First()

	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	opts := &drone.UserPatch{}
	if c.IsSet("active") {
		v := c.Bool("active")
		opts.Active = &v
	}
	if c.IsSet("admin") {
		v := c.Bool("admin")
		opts.Admin = &v
	}
	if _, err := client.UserUpdate(login, opts); err != nil {
		return err
	}
	fmt.Printf("Successfully updated user %s\n", login)
	return nil
}
