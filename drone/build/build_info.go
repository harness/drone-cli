package build

import (
	"os"
	"text/template"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/urfave/cli"
)

var buildInfoCmd = cli.Command{
	Name:      "info",
	Usage:     "show build details",
	ArgsUsage: "<repo/name> [build]",
	Action:    buildInfo,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplBuildInfo,
		},
	},
}

func buildInfo(c *cli.Context) error {
	repo := c.Args().First()
	owner, name, err := internal.ParseRepo(repo)
	if err != nil {
		return err
	}
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	_, number, err := getBuildWithArg(c.Args().Get(1), owner, repo, client)
	if err != nil {
		return err
	}

	build, err := client.Build(owner, name, number)
	if err != nil {
		return err
	}

	tmpl, err := template.New("_").Parse(c.String("format"))
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, build)
}

// template for build information
var tmplBuildInfo = `Number: {{ .Number }}
Status: {{ .Status }}
Event: {{ .Event }}
Commit: {{ .Commit }}
Branch: {{ .Branch }}
Ref: {{ .Ref }}
Message: {{ .Message }}
Author: {{ .Author }}
`
