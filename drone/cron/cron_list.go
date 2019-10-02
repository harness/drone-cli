package cron

import (
	"html/template"
	"os"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/funcmap"
	"github.com/urfave/cli"
)

var cronListCmd = cli.Command{
	Name:      "ls",
	Usage:     "list cron jobs",
	ArgsUsage: "[repo/name]",
	Action:    cronList,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "format",
			Usage:  "format output",
			Value:  tmplCronList,
			Hidden: true,
		},
	},
}

func cronList(c *cli.Context) error {
	slug := c.Args().First()
	owner, name, err := internal.ParseRepo(slug)
	if err != nil {
		return err
	}
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	list, err := client.CronList(owner, name)
	if err != nil {
		return err
	}
	format := c.String("format") + "\n"
	tmpl, err := template.New("_").Funcs(funcmap.Funcs).Parse(format)
	if err != nil {
		return err
	}
	for _, cron := range list {
		tmpl.Execute(os.Stdout, cron)
	}
	return nil
}

// template for build list information
var tmplCronList = "\x1b[33m{{ .Name }} \x1b[0m" + `
Expr: {{ .Expr }}
Next: {{ .Next | time }}
{{- if .Disabled }}
Disabled: true
{{- end }}
`
