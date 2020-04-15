package autoscale

import (
	"os"
	"text/template"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/funcmap"
	"github.com/urfave/cli"
)

var autoscaleVersionCmd = cli.Command{
	Name:   "version",
	Usage:  "server version",
	Action: autoscaleVersion,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplAutoscaleVersion,
		},
	},
}

func autoscaleVersion(c *cli.Context) error {
	client, err := internal.NewAutoscaleClient(c)
	if err != nil {
		return err
	}

	version, err := client.AutoscaleVersion()
	if err != nil {
		return err
	}

	tmpl, err := template.New("_").Funcs(funcmap.Funcs).Parse(c.String("format") + "\n")
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, version)
}

var tmplAutoscaleVersion = `Version: {{ .Version }}
Commit: {{ .Commit }}
Source: {{ .Source }}
`
