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
			Name:  "protected",
			Usage: "repository is protected",
		},
		cli.Int64Flag{
			Name:  "throttle",
			Usage: "repository throttle",
		},
		cli.DurationFlag{
			Name:  "timeout",
			Usage: "repository timeout",
		},
		cli.StringFlag{
			Name:  "visibility",
			Usage: "repository visibility",
		},
		cli.BoolFlag{
			Name:  "ignore-forks",
			Usage: "ignore forks",
		},
		cli.BoolFlag{
			Name:  "ignore-pull-requests",
			Usage: "ignore pull requests",
		},
		cli.BoolFlag{
			Name:  "auto-cancel-pull-requests",
			Usage: "automatically cancel pending pull request builds",
		},
		cli.BoolFlag{
			Name:  "auto-cancel-pushes",
			Usage: "automatically cancel pending push builds",
		},
		cli.BoolFlag{
			Name:  "auto-cancel-running",
			Usage: "automatically cancel running builds if newer commit pushed",
		},
		cli.StringFlag{
			Name:  "config",
			Usage: "repository configuration path (e.g. .drone.yml)",
		},
		cli.Int64Flag{
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
		visibility    = c.String("visibility")
		config        = c.String("config")
		timeout       = c.Duration("timeout")
		trusted       = c.Bool("trusted")
		throttle      = c.Int64("throttle")
		protected     = c.Bool("protected")
		ignoreForks   = c.Bool("ignore-forks")
		ignorePulls   = c.Bool("ignore-pull-requests")
		cancelPulls   = c.Bool("auto-cancel-pull-requests")
		cancelPush    = c.Bool("auto-cancel-pushes")
		cancelRunning = c.Bool("auto-cancel-running")
		buildCounter  = c.Int64("build-counter")
		unsafe        = c.Bool("unsafe")
	)

	patch := new(drone.RepoPatch)
	if c.IsSet("trusted") {
		patch.Trusted = &trusted
	}
	if c.IsSet("protected") {
		patch.Protected = &protected
	}
	if c.IsSet("throttle") {
		patch.Throttle = &throttle
	}
	if c.IsSet("timeout") {
		v := int64(timeout / time.Minute)
		patch.Timeout = &v
	}
	if c.IsSet("config") {
		patch.Config = &config
	}
	if c.IsSet("ignore-forks") {
		patch.IgnoreForks = &ignoreForks
	}
	if c.IsSet("ignore-pull-requests") {
		patch.IgnorePulls = &ignorePulls
	}
	if c.IsSet("auto-cancel-pull-requests") {
		patch.CancelPulls = &cancelPulls
	}
	if c.IsSet("auto-cancel-pushes") {
		patch.CancelPush = &cancelPush
	}
	if c.IsSet("auto-cancel-running") {
		patch.CancelRunning = &cancelRunning
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
		patch.Counter = &buildCounter
	}

	if _, err := client.RepoUpdate(owner, name, patch); err != nil {
		return err
	}
	fmt.Printf("Successfully updated repository %s/%s\n", owner, name)
	return nil
}
