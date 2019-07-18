package secret

import (
	"html/template"
	"os"

	"github.com/urfave/cli"

	"github.com/drone/drone-cli/drone/internal"
)

var secretListCmd = cli.Command{
	Name:      "ls",
	Usage:     "list secrets",
	ArgsUsage: "[repo/name]",
	Action:    secretList,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "repository",
			Usage: "repository name (e.g. octocat/hello-world)",
		},
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplSecretList,
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
	tmpl, err := template.New("_").Parse(format)
	if err != nil {
		return err
	}
	for _, secret := range list {
		tmpl.Execute(os.Stdout, secret)
	}
	return nil
}

// template for secret list items
var tmplSecretList = "\x1b[33m{{ .Name }} \x1b[0m" + `
Pull Request Read:  {{ .PullRequest }}
Pull Request Write: {{ .PullRequestPush }}
`
