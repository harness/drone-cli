package template

import (
	"errors"
	"html/template"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/funcmap"
	"github.com/urfave/cli"
	"os"
)

var templateInfoCmd = cli.Command{
	Name:      "info",
	Usage:     "display template info",
	ArgsUsage: "[namespace] [name]",
	Action:    templateInfo,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "namespace",
			Usage: "organization namespace",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "template name",
		},
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplTemplateInfoList,
		},
	},
}

func templateInfo(c *cli.Context) error {
	var (
		namespace    = c.String("namespace")
		templateName = c.String("name")
		format       = c.String("format") + "\n"
	)
	if templateName == "" {
		return errors.New("Missing template name")
	}
	if namespace == "" {
		return errors.New("Missing namespace")
	}

	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	templates, err := client.Template(namespace, templateName)
	if err != nil {
		return err
	}
	tmpl, err := template.New("_").Funcs(funcmap.Funcs).Parse(format)
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, templates)
}

var tmplTemplateInfoList = "\x1b[33m{{ .Name }} \x1b[0m" + "\x1b [33m{{ .Namespace }} \x1b[0m" + `
Data:  {{ .Data }}
`
