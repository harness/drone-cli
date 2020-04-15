package orgsecret

import (
	"html/template"
	"os"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/funcmap"
	"github.com/urfave/cli"
)

var secretInfoCmd = cli.Command{
	Name:      "info",
	Usage:     "display secret info",
	ArgsUsage: "[organization] [name]",
	Action:    secretInfo,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplSecretList,
		},
	},
}

func secretInfo(c *cli.Context) error {
	var (
		namespace = c.Args().First()
		name      = c.Args().Get(1)
		format    = c.String("format") + "\n"
	)
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	secret, err := client.OrgSecret(namespace, name)
	if err != nil {
		return err
	}
	tmpl, err := template.New("_").Funcs(funcmap.Funcs).Parse(format)
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, secret)
}
