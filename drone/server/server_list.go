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
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplServerList,
		},
	},
}

func serverList(c *cli.Context) error {
	client, err := internal.NewAutoscaleClient(c)
	if err != nil {
		return err
	}

	servers, err := client.ServerList()
	if err != nil || len(servers) == 0 {
		return err
	}

	tmpl, err := template.New("_").Parse(c.String("format") + "\n")
	if err != nil {
		return err
	}
	for _, server := range servers {
		tmpl.Execute(os.Stdout, server)
	}
	return nil
}

// template for server list items
var tmplServerList = `{{ .Name }}`
