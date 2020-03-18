package build

import (
	"os"
	"text/template"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/funcmap"
	"github.com/urfave/cli"
)

var buildCreateCmd = cli.Command{
	Name:      "create",
	Usage:     "create a build",
	ArgsUsage: "<repo/name>",
	Action:    buildCreate,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "commit",
			Usage: "source commit",
		},
		cli.StringFlag{
			Name:  "branch",
			Usage: "source branch",
		},
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplBuildInfo,
		},
	},
}

func buildCreate(c *cli.Context) (err error) {
	repo := c.Args().First()
	owner, name, err := internal.ParseRepo(repo)
	if err != nil {
		return err
	}

	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	build, err := client.BuildCreate(owner, name, c.String("commit"), c.String("branch"))
	if err != nil {
		return err
	}

	tmpl, err := template.New("_").Funcs(funcmap.Funcs).Parse(c.String("format"))
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, build)
}
