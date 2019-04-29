package build

import (
	"errors"
	"os"
	"strconv"
	"text/template"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/urfave/cli"
)

var buildStartCmd = cli.Command{
	Name:      "restart",
	Usage:     "restart a build",
	ArgsUsage: "<repo/name> [build]",
	Action:    buildStart,
	Flags: []cli.Flag{
		cli.StringSliceFlag{
			Name:  "param, p",
			Usage: "custom parameters to be injected into the job environment. Format: KEY=value",
		},
		cli.StringFlag{
			Name:   "format",
			Usage:  "format output",
			Value:  tmplBuildInfo,
			Hidden: true,
		},
	},
}

func buildStart(c *cli.Context) (err error) {
	repo := c.Args().First()
	owner, name, err := internal.ParseRepo(repo)
	if err != nil {
		return err
	}

	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	buildArg := c.Args().Get(1)
	var number int
	if buildArg == "last" {
		// Fetch the build number from the last build
		build, err := client.BuildLast(owner, name, "")
		if err != nil {
			return err
		}
		number = int(build.Number)
	} else {
		if len(buildArg) == 0 {
			return errors.New("missing job number")
		}
		number, err = strconv.Atoi(buildArg)
		if err != nil {
			return err
		}
	}

	params := internal.ParseKeyPair(c.StringSlice("param"))

	build, err := client.BuildRestart(owner, name, number, params)
	if err != nil {
		return err
	}

	tmpl, err := template.New("_").Parse(c.String("format"))
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, build)
}
