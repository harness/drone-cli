package build

import (
	"fmt"
	"strconv"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"
	"github.com/urfave/cli"
)

var buildCreateCmd = cli.Command{
	Name:      "trigger",
	Usage:     "trigger a build",
	ArgsUsage: "<repo/name>",
	Action:    buildCreate,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "branch, b",
			Usage: "branch",
		},
		cli.StringFlag{
			Name:  "commit, c",
			Usage: "commit sha",
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

	build, err = client.BuildCreate(owner, name, c.String("branch"), c.String("commit"))
	if err != nil {
		return err
	}

	fmt.Printf("Created build %s/%s#%d\n", owner, name, build.ID)
	return nil
}
