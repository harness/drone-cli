package repo

import (
	"fmt"

	"github.com/drone/drone-cli/drone/internal"

	"github.com/urfave/cli"
)

var repoRemoveCmd = cli.Command{
	Name:      "disable",
	Usage:     "disable a repository",
	ArgsUsage: "<repo/name>",
	Action:    repoRemove,
}

func repoRemove(c *cli.Context) error {
	repo := c.Args().First()
	owner, name, err := internal.ParseRepo(repo)
	if err != nil {
		return err
	}

	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	if err := client.RepoDisable(owner, name); err != nil {
		return err
	}
	fmt.Printf("Successfully removed repository %s/%s\n", owner, name)
	return nil
}
