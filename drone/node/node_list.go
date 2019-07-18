package node

import (
	"html/template"
	"os"

	"github.com/urfave/cli"

	"github.com/drone/drone-cli/drone/internal"
)

var nodeListCmd = cli.Command{
	Name:   "ls",
	Usage:  "list nodes",
	Action: nodeList,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplNodeList,
		},
	},
}

func nodeList(c *cli.Context) error {
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	list, err := client.NodeList()
	if err != nil {
		return err
	}
	format := c.String("format") + "\n"
	tmpl, err := template.New("_").Parse(format)
	if err != nil {
		return err
	}
	for _, cron := range list {
		tmpl.Execute(os.Stdout, cron)
	}
	return nil
}

// template for node list information
var tmplNodeList = "\x1b[33m{{ .Name }} \x1b[0m" + `
Address:  {{ .Address }}
Platform: {{ .OS }}/{{ .Arch }}
`
