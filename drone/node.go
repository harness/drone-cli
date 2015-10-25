package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/drone/drone-go/drone"
)

var MachineCmd = cli.Command{
	Name:  "node",
	Usage: "manage build nodes",
	Subcommands: []cli.Command{
		// Node List
		{
			Name:  "ls",
			Usage: "list all nodes",
			Action: func(c *cli.Context) {
				handle(c, NodeListCmd)
			},
		},
		// Node Info
		{
			Name:  "info",
			Usage: "show node details",
			Action: func(c *cli.Context) {
				handle(c, NodeInfoCmd)
			},
		},
		// Node Add
		{
			Name:  "create",
			Usage: "creates a node",
			Action: func(c *cli.Context) {
				handle(c, NodeCreateCmd)
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					EnvVar: "DOCKER_HOST",
					Name:   "docker-host",
					Usage:  "docker deamon address",
					Value:  "",
				},
				cli.BoolFlag{
					EnvVar: "DOCKER_TLS_VERIFY",
					Name:   "docker-tls-verify",
					Usage:  "docker daemon supports tlsverify",
				},
				cli.StringFlag{
					EnvVar: "DOCKER_CERT_PATH",
					Name:   "docker-cert-path",
					Usage:  "docker certificate directory",
					Value:  "",
				},
			},
		},
		// Node Delete
		{
			Name:  "rm",
			Usage: "remove a node",
			Action: func(c *cli.Context) {
				handle(c, NodeDelCmd)
			},
		},
	},
}

func NodeInfoCmd(c *cli.Context, client drone.Client) error {
	id, err := strconv.ParseInt(c.Args().Get(0), 0, 64)
	if err != nil {
		return fmt.Errorf("Invalid or missing node id. Must be an integer")
	}

	node, err := client.Node(id)
	if err != nil {
		return fmt.Errorf("Endpoint is not yet supported")
	}
	fmt.Println(node.Addr)
	return nil
}

func NodeListCmd(c *cli.Context, client drone.Client) error {
	nodes, err := client.NodeList()
	if err != nil {
		return err
	}

	for _, node := range nodes {
		fmt.Println(node.ID, node.Addr)
	}

	return nil
}

func NodeDelCmd(c *cli.Context, client drone.Client) error {
	id, err := strconv.ParseInt(c.Args().Get(0), 0, 64)
	if err != nil {
		return fmt.Errorf("Invalid or missing node id. Must be an integer")
	}

	err = client.NodeDel(id)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully removed node %d\n", id)
	return nil
}

func NodeCreateCmd(c *cli.Context, client drone.Client) error {
	node := drone.Node{
		Addr: c.String("docker-host"),
		Arch: "linux_amd64",
	}

	cert, _ := ioutil.ReadFile(filepath.Join(
		c.String("docker-cert-path"),
		"cert.pem",
	))

	key, _ := ioutil.ReadFile(filepath.Join(
		c.String("docker-cert-path"),
		"key.pem",
	))

	ca, _ := ioutil.ReadFile(filepath.Join(
		c.String("docker-cert-path"),
		"ca.pem",
	))

	if len(cert) == 0 || len(key) == 0 {
		return fmt.Errorf("Error reading cert.pem or key.pem from %s",
			c.String("docker-cert-path"))
	}

	node.Cert = string(cert)
	node.Key = string(key)

	// only use the certificate authority if tls verify
	// is enabled for this docker host.
	if c.Bool("docker-tls-verify") {
		node.CA = string(ca)
	}

	_, err := client.NodePost(&node)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully added %s\n", node.Addr)
	return nil
}
