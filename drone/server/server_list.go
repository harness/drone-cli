package server

import (
	"fmt"
	"os"
	"text/tabwriter"
	"text/template"
	"time"

	"github.com/docker/go-units"
	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"
	"github.com/drone/funcmap"
	"github.com/urfave/cli"
)

var serverListCmd = cli.Command{
	Name:      "ls",
	Usage:     "list all servers",
	ArgsUsage: " ",
	Action:    serverList,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "a, all",
			Usage: "include stopped servers",
		},
		cli.BoolFlag{
			Name:  "l, long",
			Usage: "list in long format",
		},
		cli.BoolTFlag{
			Name:  "H, headers",
			Usage: "include column headers",
		},
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplServerList,
		},
		cli.BoolFlag{
			Name:   "la",
			Hidden: true,
		},
	},
}

func serverList(c *cli.Context) error {
	client, err := internal.NewAutoscaleClient(c)
	if err != nil {
		return err
	}
	a := c.Bool("a")
	l := c.Bool("l")
	h := c.BoolT("H")

	if c.BoolT("la") {
		l = true
		a = true
	}

	servers, err := client.ServerList()
	if err != nil || len(servers) == 0 {
		return err
	}

	if l && h {
		printLong(servers, a, h)
		return nil
	}

	tmpl, err := template.New("_").Funcs(funcmap.Funcs).Parse(c.String("format") + "\n")
	if err != nil {
		return err
	}

	for _, server := range servers {
		if !a && server.State == "stopped" {
			continue
		}
		tmpl.Execute(os.Stdout, server)
	}
	return nil
}

func printLong(servers []*drone.Server, a, h bool) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	if h {
		fmt.Fprintln(w, "Name\tAddress\tState\tCreated")
	}
	for _, server := range servers {
		if !a && server.State == "stopped" {
			continue
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s ago\n",
			server.Name,
			server.Address,
			server.State,
			units.HumanDuration(
				time.Now().Sub(
					time.Unix(server.Created, 0),
				),
			),
		)
	}
	w.Flush()
}

// template for server list items
var tmplServerList = `{{ .Name }}`
