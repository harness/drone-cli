package node

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"

	"github.com/urfave/cli"
)

func getMachineHome() (path string) {
	user, err := user.Current()
	if err == nil {
		return filepath.Join(user.HomeDir, ".docker", "machine")
	}
	return
}

var nodeImportCmd = cli.Command{
	Name:   "import",
	Usage:  "import a node from docker-machine",
	Action: nodeImport,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "node name",
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

func nodeImport(c *cli.Context) error {
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	name := c.String("name")
	if name == "" {
		name = c.Args().First()
	}

	home := c.String("storage-path")
	base := filepath.Join(home, "machines", name)

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

	format := c.String("format")
	tmpl, err := template.New("_").Parse(format)
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, node)
}

type machine struct {
	DriverName string
	Driver     struct {
		IPAddress         string `json:"IPAddress"`
		MachineName       string `json:"MachineName"`
		SSHUser           string `json:"SSHUser"`
		SSHPort           int    `json:"SSHPort"`
		SSHKeyPath        string `json:"SSHKeyPath"`
		StorePath         string `json:"StorePath"`
		SwarmMaster       bool   `json:"SwarmMaster"`
		SwarmHost         string `json:"SwarmHost"`
		SwarmDiscovery    string `json:"SwarmDiscovery"`
		AccessToken       string `json:"AccessToken"`
		DropletID         int    `json:"DropletID"`
		DropletName       string `json:"DropletName"`
		Image             string `json:"Image"`
		Region            string `json:"Region"`
		SSHKeyID          int    `json:"SSHKeyID"`
		SSHKeyFingerprint string `json:"SSHKeyFingerprint"`
		SSHKey            string `json:"SSHKey"`
		Size              string `json:"Size"`
		IPv6              bool   `json:"IPv6"`
		Backups           bool   `json:"Backups"`
		PrivateNetworking bool   `json:"PrivateNetworking"`
		UserDataFile      string `json:"UserDataFile"`
		Monitoring        bool   `json:"Monitoring"`
		Tags              string `json:"Tags"`
	}
	HostOptions struct {
		AuthOptions struct {
			CertDir              string        `json:"CertDir"`
			CaCertPath           string        `json:"CaCertPath"`
			CaPrivateKeyPath     string        `json:"CaPrivateKeyPath"`
			CaCertRemotePath     string        `json:"CaCertRemotePath"`
			ServerCertPath       string        `json:"ServerCertPath"`
			ServerKeyPath        string        `json:"ServerKeyPath"`
			ClientKeyPath        string        `json:"ClientKeyPath"`
			ServerCertRemotePath string        `json:"ServerCertRemotePath"`
			ServerKeyRemotePath  string        `json:"ServerKeyRemotePath"`
			ClientCertPath       string        `json:"ClientCertPath"`
			ServerCertSANs       []interface{} `json:"ServerCertSANs"`
			StorePath            string        `json:"StorePath"`
		}
	}
}
