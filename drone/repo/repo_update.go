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
		cli.BoolFlag{
			Name:  "allow-push",
			Usage: "validate updating the build-counter is unsafe",
		},
		cli.BoolFlag{
			Name:  "allow-pr",
			Usage: "validate updating the build-counter is unsafe",
		},
		cli.BoolFlag{
			Name:  "allow-tag",
			Usage: "validate updating the build-counter is unsafe",
		},
		cli.BoolFlag{
			Name:  "allow-deploy",
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
		protected    = c.Bool("protected")
		buildCounter = c.Int("build-counter")
		unsafe       = c.Bool("unsafe")
		allowPush    = c.Bool("allow-push")
		allowPr      = c.Bool("allow-pr")
		allowTags    = c.Bool("allow-tag")
		allowDeploys = c.Bool("allow-deploy")
	)

	patch := new(drone.RepoPatch)
	if c.IsSet("trusted") {
		patch.Trusted = &trusted
	}
	if c.IsSet("protected") {
		patch.Protected = &protected
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
		patch.Counter = &buildCounter
	}
	if c.IsSet("allow-push") {
		patch.AllowPush = &allowPush
	}
	if c.IsSet("allow-pr") {
		patch.AllowPr = &allowPr
	}
	if c.IsSet("allow-tag") {
		patch.AllowTag = &allowTags
	}
	if c.IsSet("allow-deploy") {
		patch.AllowDeploy = &allowDeploys
	}

	if _, err := client.RepoUpdate(owner, name, patch); err != nil {
		return err
	}
	fmt.Printf("Successfully updated repository %s/%s\n", owner, name)
	return nil
}
