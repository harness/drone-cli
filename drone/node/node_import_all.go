package node

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"
	"github.com/drone/funcmap"
	"github.com/urfave/cli"
)

var nodeImportAllCmd = cli.Command{
	Name:   "import-all",
	Usage:  "import all node from docker-machine",
	Action: nodeImportAll,
	Flags: []cli.Flag{
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
		cli.BoolFlag{
			Name:  "paused",
			Usage: "node is paused",
		},
		cli.BoolFlag{
			Name:  "protected",
			Usage: "node is protected from deletion",
		},
		cli.StringFlag{
			Name:   "storage-path",
			Usage:  "docker machine storage path",
			Value:  getMachineHome(),
			EnvVar: "MACHINE_STORAGE_PATH",
		},
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplNodeInfo,
		},
	},
}

func nodeImportAll(c *cli.Context) error {
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	home := c.String("storage-path")

	matches, err := filepath.Glob(filepath.Join(home, "machines", "*"))
	if err != nil {
		return err
	}

	nodes, err := client.NodeList()
	if err != nil {
		return err
	}
	nodeIndex := map[string]*drone.Node{}
	for _, node := range nodes {
		nodeIndex[node.Name] = node
	}

	format := c.String("format") + "\n"
	tmpl, err := template.New("_").Funcs(funcmap.Funcs).Parse(format)
	if err != nil {
		return err
	}

	for _, name := range matches {
		base := filepath.Join(home, "machines", name)

		// if the node already exists it should be
		// ignored by the system.
		existing, ok := nodeIndex[name]
		if ok {
			tmpl.Execute(os.Stdout, existing)
			continue
		}

		conf := new(machine)
		confpath := filepath.Join(base, "config.json")
		confdata, err := ioutil.ReadFile(confpath)
		if err != nil {
			return err
		}

		err = json.Unmarshal(confdata, conf)
		if err != nil {
			return err
		}

		cakey, err := ioutil.ReadFile(conf.HostOptions.AuthOptions.CaPrivateKeyPath)
		if err != nil {
			return err
		}
		cacert, err := ioutil.ReadFile(conf.HostOptions.AuthOptions.CaCertPath)
		if err != nil {
			return err
		}
		tlskey, err := ioutil.ReadFile(conf.HostOptions.AuthOptions.ClientKeyPath)
		if err != nil {
			return err
		}
		tlscert, err := ioutil.ReadFile(conf.HostOptions.AuthOptions.ClientCertPath)
		if err != nil {
			return err
		}

		node := &drone.Node{
			UID:       fmt.Sprint(conf.Driver.DropletID),
			Provider:  c.String("provider"),
			State:     c.String("state"),
			Name:      conf.Driver.MachineName,
			Image:     conf.Driver.Image,
			Region:    conf.Driver.Region,
			Size:      conf.Driver.Size,
			OS:        c.String("os"),
			Arch:      c.String("arch"),
			Address:   conf.Driver.IPAddress,
			Capacity:  2,
			CAKey:     cakey,
			CACert:    cacert,
			TLSKey:    tlskey,
			TLSCert:   tlscert,
			Paused:    c.Bool("paused"),
			Protected: c.Bool("protected"),
		}

		_, err = client.NodeCreate(node)
		if err != nil {
			return err
		}

		tmpl.Execute(os.Stdout, node)
	}

	return nil
}
