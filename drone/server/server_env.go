package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/urfave/cli"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"
)

var serverEnvCmd = cli.Command{
	Name:      "env",
	ArgsUsage: "<servername>",
	Action:    serverEnv,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "shell",
			Usage: "specify the shell [bash, fish]",
			Value: "bash",
		},
		cli.BoolFlag{
			Name:  "clear",
			Usage: "clear cert cache",
		},
	},
}

func serverEnv(c *cli.Context) error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	name := c.Args().First()
	if len(name) == 0 {
		return fmt.Errorf("Missing or invalid server name")
	}

	home := path.Join(u.HomeDir, ".drone", "certs")
	base := path.Join(home, name)

	if c.Bool("clean") {
		os.RemoveAll(home)
	}

	server := new(drone.Server)
	if _, err := os.Stat(base); err == nil {
		data, err := ioutil.ReadFile(path.Join(base, "server.json"))
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, server)
		if err != nil {
			return err
		}
	} else {
		client, err := internal.NewAutoscaleClient(c)
		if err != nil {
			return err
		}
		server, err = client.Server(name)
		if err != nil {
			return err
		}
		data, err := json.Marshal(server)
		if err != nil {
			return err
		}
		err = os.MkdirAll(base, 0755)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path.Join(base, "server.json"), data, 0644)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path.Join(base, "ca.pem"), server.CACert, 0644)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path.Join(base, "cert.pem"), server.TLSCert, 0644)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path.Join(base, "key.pem"), server.TLSKey, 0644)
		if err != nil {
			return err
		}
	}

	switch c.String("shell") {
	case "fish":
		fmt.Fprintf(os.Stdout, fishf, base, server.Address, server.Name)
	case "powershell":
		fmt.Fprintf(os.Stdout, powershellf, base, server.Address, server.Name)
	default:
		fmt.Fprintf(os.Stdout, bashf, base, server.Address, server.Name)
	}

	return nil
}

var bashf = `
export DOCKER_TLS=1
export DOCKER_TLS_VERIFY=
export DOCKER_CERT_PATH=%q
export DOCKER_HOST=tcp://%s:2376

# Run this command to configure your shell:
# eval "$(drone server env %s)"
`

var fishf = `
sex -x DOCKER_TLS "1";
set -x DOCKER_TLS_VERIFY "";
set -x DOCKER_CERT_PATH %q;
set -x DOCKER_HOST tcp://%s:2376;

# Run this command to configure your shell:
# eval "$(drone server env %s --shell=fish)"
`

var powershellf = `
$Env:DOCKER_TLS = "1"
$Env:DOCKER_TLS_VERIFY = ""
$Env:DOCKER_CERT_PATH = %q
$Env:DOCKER_HOST = "tcp://%s:2376"

# Run this command to configure your shell:
# drone server env %s --shell=powershell | Invoke-Expression
`
