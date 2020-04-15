package server

import (
	"fmt"
	"os"
	"text/template"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/funcmap"
	"github.com/urfave/cli"
)

var serverDestroyCmd = cli.Command{
	Name:      "destroy",
	Usage:     "destroy a server",
	ArgsUsage: "<servername>",
	Action:    serverDestroy,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplServerDestroy,
		},
		cli.BoolFlag{
			Name:  "force",
			Usage: "force destroy",
		},
	},
}

func serverDestroy(c *cli.Context) error {
	client, err := internal.NewAutoscaleClient(c)
	if err != nil {
		return err
	}

	name := c.Args().First()
	if len(name) == 0 {
		return fmt.Errorf("Missing or invalid server name")
	}

	err = client.ServerDelete(name, c.Bool("force"))
	if err != nil {
		return err
	}

	server, err := client.Server(name)
	if err != nil {
		return err
	}

	tmpl, err := template.New("_").Funcs(funcmap.Funcs).Parse(c.String("format") + "\n")
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, server)
}

var tmplServerDestroy = `Name: {{ .Name }}
Address: {{ .Address }}
Region: {{ .Region }}
Size: {{.Size}}
State: {{ .State }}
`
