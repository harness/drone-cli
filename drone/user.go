package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/codegangsta/cli"
	"github.com/drone/drone-go/drone"
)

var userCmd = cli.Command{
	Name:  "user",
	Usage: "manage users",
	Subcommands: []cli.Command{

		// list command
		cli.Command{
			Name:   "ls",
			Usage:  "list all users",
			Action: userList,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "format",
					Usage: "format output",
					Value: tmplUserList,
				},
			},
		},

		// info command
		cli.Command{
			Name:   "info",
			Usage:  "show user details",
			Action: userInfo,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "format",
					Usage: "format output",
					Value: tmplUserInfo,
				},
			},
		},

		// add command
		cli.Command{
			Name:   "add",
			Usage:  "adds a user",
			Action: userAdd,
		},

		// remove command
		cli.Command{
			Name:   "rm",
			Usage:  "remove a user",
			Action: userRemove,
		},
	},
}

// command to de-register a user from drone server.
func userRemove(c *cli.Context) error {
	login := c.Args().First()

	client, err := newClient(c)
	if err != nil {
		return err
	}
	if err := client.UserDel(login); err != nil {
		return err
	}
	fmt.Printf("Successfully removed user %s\n", login)
	return nil
}

// command to register a user with drone server.
func userAdd(c *cli.Context) error {
	login := c.Args().First()

	client, err := newClient(c)
	if err != nil {
		return err
	}
	user, err := client.UserPost(&drone.User{Login: login})
	if err != nil {
		return err
	}
	fmt.Printf("Successfully added user %s\n", user.Login)
	return nil
}

// command to display information for a registered user.
func userInfo(c *cli.Context) error {
	client, err := newClient(c)
	if err != nil {
		return err
	}

	login := c.Args().First()
	if len(login) == 0 {
		return fmt.Errorf("Missing or invalid user login")
	}
	user, err := client.User(login)
	if err != nil {
		return err
	}

	tmpl, err := template.New("_").Parse(c.String("format") + "\n")
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, user)
}

// command to list registered user.
func userList(c *cli.Context) error {
	client, err := newClient(c)
	if err != nil {
		return err
	}

	users, err := client.UserList()
	if err != nil || len(users) == 0 {
		return err
	}

	tmpl, err := template.New("_").Parse(c.String("format") + "\n")
	if err != nil {
		return err
	}
	for _, user := range users {
		tmpl.Execute(os.Stdout, user)
	}
	return nil
}

// template for user list.
var tmplUserList = `{{ .Login }}`

// template for user info
var tmplUserInfo = `User: {{ .Login }}
Email: {{ .Email }}`
