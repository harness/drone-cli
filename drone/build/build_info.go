package build

import (
	"os"
	"text/template"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/urfave/cli"
	"strconv"
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

	buildArgStr := c.Args().Get(1)
	var number int
	if buildArgStr == "last" {
		build, err := client.BuildLast(owner, name, "")
		if err != nil {
			return err
		}
		number = build.Number
	} else if parsedNumber, err := strconv.Atoi(buildArgStr); err != nil {
		return errInvalidBuildNumber
	} else {
		number = parsedNumber
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
