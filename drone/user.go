package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/drone/drone-go/drone"
)

var UserCmd = cli.Command{
	Name:  "user",
	Usage: "manage users",
	Subcommands: []cli.Command{
		// User List
		{
			Name:  "ls",
			Usage: "list all users",
			Action: func(c *cli.Context) {
				handle(c, UserListCmd)
			},
		},
		// User Info
		{
			Name:  "info",
			Usage: "show user details",
			Action: func(c *cli.Context) {
				handle(c, UserInfoCmd)
			},
		},
		// User Add
		{
			Name:  "add",
			Usage: "adds a user",
			Action: func(c *cli.Context) {
				handle(c, UserAddCmd)
			},
		},
		// User Delete
		{
			Name:  "rm",
			Usage: "remove a user",
			Action: func(c *cli.Context) {
				handle(c, UserDelCmd)
			},
		},
		// User Self
		{
			Name:  "self",
			Usage: "show the current user details",
			Flags: []cli.Flag{},
			Action: func(c *cli.Context) {
				handle(c, UserSelfCmd)
			},
		},
	},
}

func UserInfoCmd(c *cli.Context, client drone.Client) error {
	login := c.Args().Get(0)
	if len(login) == 0 {
		return fmt.Errorf("Missing or invalid user login")
	}

	user, err := client.User(login)
	if err != nil {
		return err
	}
	fmt.Println(user.Login)
	fmt.Println(user.Email)
	fmt.Println(user.Admin)
	fmt.Println(user.Active)
	return nil
}

func UserListCmd(c *cli.Context, client drone.Client) error {

	users, err := client.UserList()
	if err != nil {
		return err
	}
	for _, user := range users {
		fmt.Println(user.Login)
		fmt.Println("\t", user.Email)
		fmt.Println("\t", user.Admin)
		fmt.Println("\t", user.Active)
	}

	return nil
}

func UserAddCmd(c *cli.Context, client drone.Client) error {
	login := c.Args().Get(0)
	if len(login) == 0 {
		return fmt.Errorf("Missing or invalid user login")
	}

	user, err := client.UserPost(&drone.User{Login: login})
	if err != nil {
		return err
	}
	fmt.Printf("Successfully added user %s\n", user.Login)
	return nil
}

func UserDelCmd(c *cli.Context, client drone.Client) error {
	login := c.Args().Get(0)
	if len(login) == 0 {
		return fmt.Errorf("Missing or invalid user login")
	}

	err := client.UserDel(login)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully removed user %s\n", login)
	return nil
}

func UserSelfCmd(c *cli.Context, client drone.Client) error {
	user, err := client.Self()
	if err != nil {
		return err
	}

	fmt.Println(user.Login)
	return nil
}
