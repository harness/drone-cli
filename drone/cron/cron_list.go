package cron

import (
	"html/template"
	"os"
	"time"

	"github.com/urfave/cli"

	"github.com/drone/drone-cli/drone/internal"
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
	tmpl, err := template.New("_").Funcs(funcs).Parse(format)
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
Next: {{ fromUnix .Next }}
{{- if .Disabled }}
Disabled: true
{{- end }}
`

var funcs = map[string]interface{}{
	"fromUnix": func(v interface{}) string {
		i, ok := v.(int64)
		if !ok {
			return ""
		}
		return time.Unix(i, 0).String()
	},
}
