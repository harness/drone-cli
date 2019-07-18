package node

import (
	"html/template"
	"os"

	"github.com/drone/drone-cli/drone/internal"

	"github.com/urfave/cli"
)

var nodeInfoCmd = cli.Command{
	Name:   "info",
	Usage:  "display node info",
	Action: nodeInfo,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplNodeInfo,
		},
	},
}

func nodeInfo(c *cli.Context) error {
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	name := c.Args().First()
	node, err := client.Node(name)
	if err != nil {
		return err
	}
	format := c.String("format")
	tmpl, err := template.New("_").Parse(format)
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, node)
}

var tmplNodeInfo = "\x1b[33m{{ .Name }} \x1b[0m" + `
Address:  {{ .Address }}
Region:   {{ .Region }}
Instance: {{ .Size }}
OS:       {{ .OS }}
Arch:     {{ .Arch }}
Locked:   {{ .Protected }}
Paused:   {{ .Paused }}
`
