package template

import (
	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/funcmap"
	"github.com/urfave/cli"
	"os"
	"text/template"
)

var templateListCmd = cli.Command{
	Name:      "ls",
	Usage:     "list templates",
	ArgsUsage: "[]",
	Action:    templateList,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplTemplateInfoList,
		},
	},
}

func templateList(c *cli.Context) error {
	var (
		namespace = c.Args().First()
		format    = c.String("format") + "\n"
	)
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	list, err := client.TemplateList(namespace)
	if err != nil {
		return err
	}
	tmpl, err := template.New("_").Funcs(funcmap.Funcs).Parse(format)
	if err != nil {
		return err
	}
	for _, templates := range list {
		tmpl.Execute(os.Stdout, templates)
	}
	return nil
}
