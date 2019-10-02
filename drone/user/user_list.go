package user

import (
	"os"
	"text/template"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/funcmap"
	"github.com/urfave/cli"
)

var userListCmd = cli.Command{
	Name:      "ls",
	Usage:     "list all users",
	ArgsUsage: " ",
	Action:    userList,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplUserList,
		},
	},
}

func userList(c *cli.Context) error {
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	users, err := client.UserList()
	if err != nil || len(users) == 0 {
		return err
	}

	tmpl, err := template.New("_").Funcs(funcmap.Funcs).Parse(c.String("format") + "\n")
	if err != nil {
		return err
	}
	for _, user := range users {
		tmpl.Execute(os.Stdout, user)
	}
	return nil
}

// template for user list items
var tmplUserList = `{{ .Login }}`
