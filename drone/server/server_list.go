package server

import (
	"os"
	"text/template"

	"github.com/urfave/cli"

	"github.com/drone/drone-cli/drone/internal"
)

var serverListCmd = cli.Command{
	Name:      "ls",
	Usage:     "list all servers",
	ArgsUsage: " ",
	Action:    serverList,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "a",
			Usage: "show all servers",
		},
		cli.StringFlag{
			Name:   "format",
			Usage:  "format output",
			Value:  tmplServerList,
			Hidden: true,
		},
	},
}

func serverList(c *cli.Context) error {
	client, err := internal.NewAutoscaleClient(c)
	if err != nil {
		return err
	}
	all := c.Bool("a")

	servers, err := client.ServerList()
	if err != nil || len(servers) == 0 {
		return err
	}

	tmpl, err := template.New("_").Parse(c.String("format") + "\n")
	if err != nil {
		return err
	}
	for _, server := range servers {
		if !all && server.State == "stopped" {
			continue
		}
		tmpl.Execute(os.Stdout, server)
	}
	return nil
}

// template for server list items
var tmplServerList = `{{ .Name }}`
