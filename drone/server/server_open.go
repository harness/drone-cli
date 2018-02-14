package server

import (
	"fmt"
	"net/url"

	"github.com/pkg/browser"
	"github.com/urfave/cli"

	"github.com/drone/drone-cli/drone/internal"
)

var serverOpenCmd = cli.Command{
	Name:      "open",
	Usage:     "open server dashboard",
	ArgsUsage: "<servername>",
	Action:    serverOpen,
}

func serverOpen(c *cli.Context) error {
	client, err := internal.NewAutoscaleClient(c)
	if err != nil {
		return err
	}

	name := c.Args().First()
	if len(name) == 0 {
		return fmt.Errorf("Missing or invalid server name")
	}

	server, err := client.Server(name)
	if err != nil {
		return err
	}

	uri := new(url.URL)
	uri.Scheme = "http"
	uri.Host = server.Address + ":8080"
	uri.User = url.UserPassword("admin", server.Secret)

	return browser.OpenURL(uri.String())
}
