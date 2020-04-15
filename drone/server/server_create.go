package server

import (
	"os"
	"text/template"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/funcmap"
	"github.com/urfave/cli"
)

var serverCreateCmd = cli.Command{
	Name:   "create",
	Usage:  "create a new server",
	Action: serverCreate,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplServerCreate,
		},
	},
}

func serverCreate(c *cli.Context) error {
	client, err := internal.NewAutoscaleClient(c)
	if err != nil {
		return err
	}

	server, err := client.ServerCreate()
	if err != nil {
		return err
	}

	tmpl, err := template.New("_").Funcs(funcmap.Funcs).Parse(c.String("format") + "\n")
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, server)
}

var tmplServerCreate = `Name: {{ .Name }}
State: {{ .State }}
`
