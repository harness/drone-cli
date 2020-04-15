package user

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"
)

var userBlockCmd = cli.Command{
	Name:      "block",
	Usage:     "block a user",
	ArgsUsage: "<username>",
	Action:    userBlock,
}

func userBlock(c *cli.Context) error {
	login := c.Args().First()

	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	active := false
	opts := &drone.UserPatch{
		Active: &active,
	}
	_, err = client.UserUpdate(login, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully blocked user %s\n", login)
	return nil
}
