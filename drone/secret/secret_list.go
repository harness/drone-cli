package secret

import (
	"html/template"
	"os"
	"strings"

	"github.com/urfave/cli"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"
)

var secretListCmd = cli.Command{
	Name:      "ls",
	Aliases:   []string{"list"},
	Usage:     "list secrets",
	ArgsUsage: "[repo/name]",
	Action:    secretList,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "repository",
			Usage: "repository name (e.g. octocat/hello-world)",
		},
		cli.StringFlag{
			Name:   "format",
			Usage:  "format output",
			Value:  tmplSecretList,
			Hidden: true,
		},
	},
}

func secretList(c *cli.Context) error {
	var (
		format   = c.String("format") + "\n"
		reponame = c.String("repository")
	)
	if reponame == "" {
		reponame = c.Args().First()
	}
	owner, name, err := internal.ParseRepo(reponame)
	if err != nil {
		return err
	}
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	list, err := client.SecretList(owner, name)
	if err != nil {
		return err
	}
	for _, secret := range list {
		err = printSecret(secret, format)
		if err != nil {
			return err
		}
	}
	return nil
}

func printSecret(secret *drone.Secret, format string) error {
	tmpl, err := template.New("_").Funcs(secretFuncMap).Parse(format)
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, secret)
}

// template for secret list items
var tmplSecretList = "\x1b[33m{{ .Name }} \x1b[0m" + `
Events: {{ list .Events }}
{{- if .Images }}
Images: {{ list .Images }}
{{- else }}
Images: <any>
{{- end }}
`

var secretFuncMap = template.FuncMap{
	"list": func(s []string) string {
		return strings.Join(s, ", ")
	},
}
