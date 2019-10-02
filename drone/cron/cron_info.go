package cron

import (
	"html/template"
	"os"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/funcmap"
	"github.com/urfave/cli"
)

var cronInfoCmd = cli.Command{
	Name:      "info",
	Usage:     "display cron info",
	ArgsUsage: "[repo/name]",
	Action:    cronInfo,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplCronList,
		},
	},
}

func cronInfo(c *cli.Context) error {
	slug := c.Args().First()
	owner, name, err := internal.ParseRepo(slug)
	if err != nil {
		return err
	}
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	cronjob := c.Args().Get(1)
	cron, err := client.Cron(owner, name, cronjob)
	if err != nil {
		return err
	}
	format := c.String("format")
	tmpl, err := template.New("_").Funcs(funcmap.Funcs).Parse(format)
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, cron)
}
