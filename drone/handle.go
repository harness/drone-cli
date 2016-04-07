package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/drone/drone-go/drone"
	"github.com/jackspirou/syscerts"
)

type handlerFunc func(*cli.Context, drone.Client) error

// handle wraps the command function handlers and
// sets up the environment.
func handle(c *cli.Context, fn handlerFunc) {
	var token = c.GlobalString("token")
	var server = strings.TrimSuffix(c.GlobalString("server"), "/")

	// if no server url is provided we can default
	// to the hosted Drone service.
	if len(server) == 0 {
		fmt.Println("Error: you must provide the Drone server address.")
		os.Exit(1)
	}

	if len(token) == 0 {
		fmt.Println("Error: you must provide your Drone access token.")
		os.Exit(1)
	}

	// attempt to find system CA certs
	certs := syscerts.SystemRootsPool()
	tlsConfig := &tls.Config{RootCAs: certs}

	// create the drone client with TLS options
	client := drone.NewClientTokenTLS(server, token, tlsConfig)

	// handle the function
	if err := fn(c, client); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
