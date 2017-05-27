package build

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/urfave/cli"
)

var buildStartCmd = cli.Command{
	Name:      "start",
	Usage:     "start a build",
	ArgsUsage: "<repo/name> [build]",
	Action:    buildStart,
	Flags: []cli.Flag{
		cli.StringSliceFlag{
			Name:  "param, p",
			Usage: "custom parameters to be injected into the job environment. Format: KEY=value",
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

	_, number, err := getBuildWithArg(c.Args().Get(1), owner, name, client)
	params := internal.ParseKeyPair(c.StringSlice("param"))

	build, err := client.BuildStart(owner, name, number, params)
	if err != nil {
		return err
	}

	fmt.Printf("Starting build %s/%s#%d\n", owner, name, build.Number)
	return nil
}
