package repo

import (
	"fmt"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"

	"github.com/urfave/cli"
)
var repoMoveCmd = cli.Command{
	Name:      "mv",
	Usage:     "move the repository to a new location",
	ArgsUsage: "<repo/name> <newRepo/newName>",
	Action:    repoMove,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "unsafe",
			Usage: "validate updating the full repo name is unsafe - build status in repos will reflect old name",
		},
	},
}

func repoMove(c *cli.Context) error {
	repo := c.Args().Get(0)
	newRepo := c.Args().Get(1)

	owner, name, err := internal.ParseRepo(repo)
	if err != nil {
		return err
	}


	newOwner, newName, newErr := internal.ParseRepo(newRepo)
	if newErr != nil {
		return newErr
	}

	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	var (
		unsafe = c.Bool("unsafe")
	)
	patch := new(drone.RepoPatch)
	if  !unsafe {
		fmt.Printf("Setting the build counter is an unsafe operation that could put your repository in an inconsistent state. Please use --unsafe to proceed")
		return nil
	}
	patch.Owner = &newOwner
	patch.Name = &newName
	if _, err := client.RepoPatch(owner, name, patch); err != nil {
		return err
	}
	if repairErr := client.RepoRepair(newOwner, newName); repairErr != nil {
		return repairErr
	}
	fmt.Printf("Successfully moved repository from %s/%s to %s/%s\n", owner, name, newOwner, newName)
	return nil
}