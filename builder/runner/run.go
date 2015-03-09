package main

import (
	"os"

	"github.com/drone/drone-cli/builder"
	"github.com/drone/drone-cli/common"
	"github.com/drone/drone-cli/common/config"
	"github.com/drone/drone-cli/engine/docker"
	//"github.com/drone/drone-cli/cluster"
	//"github.com/drone/drone-cli/cluster/docker"

	"github.com/samalba/dockerclient"
)

var yaml = `
build:
  image: golang:$$go_version
  commands:
    - ls -la /drone/src/github.com/drone/drone
    - go version

matrix:
  go_version:
    - 1.3.3
#   - 1.4.2
`

func main() {

	matrix, err := config.ParseMatrix(yaml)
	if err != nil {
		println(err.Error())
		return
	}

	client, err := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
	if err != nil {
		println(err.Error())
		return
	}
	engine := docker.New(client)
	if err := engine.Setup(); err != nil {
		println(err.Error())
		return
	}
	defer engine.Teardown()

	clone := &common.Clone{}
	clone.Dir = "/drone/src/github.com/drone/drone"
	clone.Branch = "master"
	clone.Sha = "4fbcc1dd41c5e2792c034d31e350f521890ad723"
	clone.Remote = "git://github.com/drone/drone.git"

	for _, config := range matrix {
		r := &builder.Request{}
		r.Clone = clone
		r.Config = config

		rw := &builder.Response{}
		rw.Writer = os.Stdout

		c := builder.Context{engine, r, rw}
		run(&c)
	}

}

func run(c *builder.Context) {
	build := builder.Builder{}
	//after := builder.Builder{}
	//always := builder.Builder{}
	defer func() {
		//always.Cancel()
		//after.Cancel()
		build.Cancel()
	}()

	// Mandatory Setup step to bootstrap our environment
	// and generate our build script
	c.Request.Config.Setup = common.Step{}
	c.Request.Config.Setup.Image = "plugin/drone-build"
	c.Request.Config.Setup.Config = c.Request.Config.Build.Config

	// Mandatory change how the build behaves
	c.Request.Config.Build.Entrypoint = []string{"/bin/bash"}
	c.Request.Config.Build.Command = []string{"/drone/bin/build.sh"}

	// loop through compose containers and setup
	// build handlers
	for _, step := range c.Request.Config.Compose {
		build.Register(builder.ServiceHandler(&step))
	}

	// setup the build and clone containers
	build.Register(builder.BatchHandler(&c.Request.Config.Setup))
	build.Register(builder.BatchHandler(&c.Request.Config.Clone))
	build.Register(builder.BatchHandler(&c.Request.Config.Build))
	err := build.Build(c.Response)
	if err != nil {
		println(err.Error())
	}

	// if c.ExitCode() == 0 {
	// 	after.Register(builder.PluginHandler(nil))
	// 	after.Register(builder.PluginHandler(nil))
	// 	after.Register(builder.PluginHandler(nil))
	// 	after.Register(builder.PluginHandler(nil))
	// 	after.Register(builder.PluginHandler(nil))
	// 	after.Build(c)
	// }

	// always.Register(builder.PluginHandler(nil))
	// always.Register(builder.PluginHandler(nil))
	// always.Register(builder.PluginHandler(nil))
	// always.Register(builder.PluginHandler(nil))
	// always.Register(builder.PluginHandler(nil))
	// always.Build(c)
}
