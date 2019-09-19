package admit

import (
	"context"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/admission"
	"github.com/urfave/cli"
)

// Command exports the admission command set.
var Command = cli.Command{
	Name:      "admit",
	Usage:     "test user admission",
	ArgsUsage: "user",
	Action:    admit,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "user",
			Usage: "username",
		},

		cli.StringFlag{
			Name:   "endpoint",
			Usage:  "plugin endpoint",
			EnvVar: "DRONE_ADMISSION_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "secret",
			Usage:  "plugin secret",
			EnvVar: "DRONE_ADMISSION_SECRET",
		},
		cli.StringFlag{
			Name:   "ssl-skip-verify",
			Usage:  "plugin ssl verification disabled",
			EnvVar: "DRONE_ADMISSION_SKIP_VERIFY",
		},
	},
}

func admit(c *cli.Context) error {
	login := c.String("user")
	if login == "" {
		login = c.Args().First()
	}

	req := &admission.Request{
		User: drone.User{
			Login: login,
		},
	}

	client := admission.Client(
		c.String("endpoint"),
		c.String("secret"),
		c.Bool("ssl-skip-verify"),
	)
	_, err := client.Admit(context.Background(), req)
	if err != nil {
		return err
	}

	return nil
}
