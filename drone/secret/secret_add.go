package secret

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"

	"github.com/urfave/cli"
)

var flags = []cli.Flag{
	cli.StringFlag{
		Name:  "repository",
		Usage: "repository name (e.g. octocat/hello-world)",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "secret name",
	},
	cli.StringFlag{
		Name:  "value",
		Usage: "secret value (@path/file.txt to read contents)",
	},
	cli.StringSliceFlag{
		Name:  "image",
		Usage: "limit to these images (repeat flag for each image)",
	},
	cli.BoolFlag{
		Name:  "push",
		Usage: "limit to push events",
	},
	cli.BoolFlag{
		Name:  "tag",
		Usage: "limit to tag events",
	},
	cli.BoolFlag{
		Name:  "deploy",
		Usage: "limit to deployment events",
	},
	cli.BoolFlag{
		Name:  "pr",
		Usage: "limit to pull_request events",
	},
	cli.BoolFlag{
		Name:  "all",
		Usage: "don't limit events",
	},
	cli.StringSliceFlag{
		Name:  "event",
		Usage: "limit to these events (repeat flag for each event) [deprecated]",
	},
}

var secretCreateCmd = cli.Command{
	Name:      "add",
	Usage:     "adds a secret",
	ArgsUsage: "[repo/name]",
	Action:    secretCreate,
	Flags:     flags,
}

func secretCreate(c *cli.Context) error {
	reponame := c.String("repository")
	if reponame == "" {
		reponame = c.Args().First()
	}
	owner, name, err := internal.ParseRepo(reponame)
	if err != nil {
		return err
	}
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	secret, err := makeSecret(c)
	if err != nil {
		return err
	}
	if len(secret.Events) == 0 {
		secret.Events = defaultSecretEvents
	}
	_, err = client.SecretCreate(owner, name, secret)
	if err == nil {
		return printSecret(secret, tmplSecretList)
	}
	return err
}

func makeSecret(c *cli.Context) (*drone.Secret, error) {
	events := c.StringSlice("event")
	for _, e := range events {
		if !validSecretEvents[e] {
			return nil, fmt.Errorf("Error: Invalid event (%v).", e)
		}
	}

	if c.Bool("push") || c.Bool("all") {
		events = append(events, drone.EventPush)
	}
	if c.Bool("tag") || c.Bool("all") {
		events = append(events, drone.EventTag)
	}
	if c.Bool("deploy") || c.Bool("all") {
		events = append(events, drone.EventDeploy)
	}
	if c.Bool("pr") || c.Bool("all") {
		events = append(events, drone.EventPull)
	}

	secret := &drone.Secret{
		Name:   c.String("name"),
		Value:  c.String("value"),
		Images: c.StringSlice("image"),
		Events: events,
	}
	if strings.HasPrefix(secret.Value, "@") {
		path := strings.TrimPrefix(secret.Value, "@")
		out, ferr := ioutil.ReadFile(path)
		if ferr != nil {
			return nil, ferr
		}
		secret.Value = string(out)
	}
	return secret, nil
}

var defaultSecretEvents = []string{
	drone.EventPush,
	drone.EventTag,
	drone.EventDeploy,
}

var validSecretEvents = map[string]bool{
	drone.EventPush:   true,
	drone.EventPull:   true,
	drone.EventTag:    true,
	drone.EventDeploy: true,
}
