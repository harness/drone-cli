package repo

import (
	"fmt"
	"time"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"

	"github.com/urfave/cli"
)

var repoUpdateCmd = cli.Command{
	Name:      "update",
	Usage:     "update a repository",
	ArgsUsage: "<repo/name>",
	Action:    repoUpdate,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "trusted",
			Usage: "repository is trusted",
		},
		cli.BoolFlag{
			Name:  "gated",
			Usage: "repository is gated",
		},
		cli.BoolFlag{
			Name:  "tag",
			Usage: "repository support tag hook",
		},
		cli.BoolFlag{
			Name:  "pull",
			Usage: "repository support pull hook",
		},
		cli.BoolFlag{
			Name:  "push",
			Usage: "repository support push hook",
		},
		cli.BoolFlag{
			Name:  "deploy",
			Usage: "repository support deploy hook",
		},
		cli.DurationFlag{
			Name:  "timeout",
			Usage: "repository timeout",
		},
		cli.StringFlag{
			Name:  "visibility",
			Usage: "repository visibility",
		},
		cli.StringFlag{
			Name:  "config",
			Usage: "repository configuration path (e.g. .drone.yml)",
		},
		cli.IntFlag{
			Name:  "build-counter",
			Usage: "repository starting build number",
		},
		cli.BoolFlag{
			Name:  "unsafe",
			Usage: "validate updating the build-counter is unsafe",
		},
	},
}

func repoUpdate(c *cli.Context) error {
	repo := c.Args().First()
	owner, name, err := internal.ParseRepo(repo)
	if err != nil {
		return err
	}

	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	var (
		visibility   = c.String("visibility")
		config       = c.String("config")
		timeout      = c.Duration("timeout")
		trusted      = c.Bool("trusted")
		gated        = c.Bool("gated")
		buildCounter = c.Int("build-counter")
		unsafe       = c.Bool("unsafe")
		tag          = c.Bool("tag")
		push         = c.Bool("push")
		pull         = c.Bool("pull")
		deploy       = c.Bool("deploy")
	)

	patch := new(drone.RepoPatch)
	if c.IsSet("tag") {
		patch.AllowTag = &tag
	}
	if c.IsSet("push") {
		patch.AllowPush = &push
	}
	if c.IsSet("pull") {
		patch.AllowPull = &pull
	}
	if c.IsSet("deploy") {
		patch.AllowDeploy = &deploy
	}
	if c.IsSet("trusted") {
		patch.IsTrusted = &trusted
	}
	if c.IsSet("gated") {
		patch.IsGated = &gated
	}
	if c.IsSet("timeout") {
		v := int64(timeout / time.Minute)
		patch.Timeout = &v
	}
	if c.IsSet("config") {
		patch.Config = &config
	}
	if c.IsSet("visibility") {
		switch visibility {
		case "public", "private", "internal":
			patch.Visibility = &visibility
		}
	}
	if c.IsSet("build-counter") && !unsafe {
		fmt.Printf("Setting the build counter is an unsafe operation that could put your repository in an inconsistent state. Please use --unsafe to proceed")
	}
	if c.IsSet("build-counter") && unsafe {
		patch.BuildCounter = &buildCounter
	}

	if _, err := client.RepoPatch(owner, name, patch); err != nil {
		return err
	}
	fmt.Printf("Successfully updated repository %s/%s\n", owner, name)
	return nil
}
