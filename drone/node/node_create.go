package node

import (
	"html/template"
	"io/ioutil"
	"os"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"

	"github.com/urfave/cli"
)

var nodeCreateCmd = cli.Command{
	Name:   "add",
	Usage:  "adds a node",
	Action: nodeCreate,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "node name",
		},
		cli.StringFlag{
			Name:  "hostname",
			Usage: "node hostname or ip address",
		},
		cli.StringFlag{
			Name:  "ca-key",
			Usage: "path to ca key",
		},
		cli.StringFlag{
			Name:  "ca-cert",
			Usage: "path to ca cert",
		},
		cli.StringFlag{
			Name:  "tls-key",
			Usage: "path to tls key",
		},
		cli.StringFlag{
			Name:  "tls-cert",
			Usage: "path to tls cert",
		},
		cli.StringFlag{
			Name:  "tls-server-name",
			Usage: "tls server name",
		},
		cli.IntFlag{
			Name:  "capacity",
			Usage: "node capacity",
			Value: 2,
		},
		cli.StringFlag{
			Name:  "os",
			Usage: "node os",
			Value: "linux",
		},
		cli.StringFlag{
			Name:  "arch",
			Usage: "node arch",
			Value: "amd64",
		},
		cli.StringFlag{
			Name:  "region",
			Usage: "node region",
		},
		cli.StringFlag{
			Name:  "instance",
			Usage: "node instance type",
		},
		cli.StringFlag{
			Name:  "image",
			Usage: "node image (i.e. ami)",
		},
		cli.StringFlag{
			Name:  "provider",
			Usage: "node hosting provider (e.g. amazon)",
		},
		cli.BoolFlag{
			Name:  "paused",
			Usage: "node is paused",
		},
		cli.BoolFlag{
			Name:  "protected",
			Usage: "node is protected from deletion",
		},
		cli.StringFlag{
			Name:   "format",
			Usage:  "format output",
			Value:  tmplNodeInfo,
			Hidden: true,
		},
	},
}

func nodeCreate(c *cli.Context) error {
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	cakey, err := ioutil.ReadFile(c.String("ca-key"))
	if err != nil {
		return err
	}
	cacert, err := ioutil.ReadFile(c.String("ca-cert"))
	if err != nil {
		return err
	}
	tlskey, err := ioutil.ReadFile(c.String("tls-key"))
	if err != nil {
		return err
	}
	tlscert, err := ioutil.ReadFile(c.String("tls-cert"))
	if err != nil {
		return err
	}

	node := &drone.Node{
		UID:       c.String("id"),
		Provider:  c.String("provider"),
		State:     c.String("state"),
		Name:      c.String("name"),
		Image:     c.String("image"),
		Region:    c.String("region"),
		Size:      c.String("instance"),
		OS:        c.String("os"),
		Arch:      c.String("arch"),
		Address:   c.String("hostname"),
		Capacity:  c.Int("capacity"),
		CAKey:     cakey,
		CACert:    cacert,
		TLSKey:    tlskey,
		TLSCert:   tlscert,
		TLSName:   c.String("tls-server-name"),
		Paused:    c.Bool("paused"),
		Protected: c.Bool("protected"),
	}

	_, err = client.NodeCreate(node)
	if err != nil {
		return err
	}

	format := c.String("format")
	tmpl, err := template.New("_").Parse(format)
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, node)
}
