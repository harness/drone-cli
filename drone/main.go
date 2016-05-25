package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	app := cli.NewApp()
	app.Name = "drone"
	app.Version = "0.5"
	app.Usage = "command line utility"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "t, token",
			Usage:  "server auth token",
			EnvVar: "DRONE_TOKEN",
		},
		cli.StringFlag{
			Name:   "s, server",
			Usage:  "server location",
			EnvVar: "DRONE_SERVER",
		},
	}
	app.Commands = []cli.Command{
		agentCmd,
		buildCmd,
		deployCmd,
		execCmd,
		infoCmd,
		secretCmd,
		signCmd,
		repoCmd,
		userCmd,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
