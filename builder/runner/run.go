package main

import (
	"os"

	"github.com/drone/drone-cli/builder"
	"github.com/drone/drone-cli/common"
	"github.com/drone/drone-cli/common/config"
	"github.com/drone/drone-cli/engine"
	"github.com/drone/drone-cli/engine/docker"
	//"github.com/drone/drone-cli/cluster"
	//"github.com/drone/drone-cli/cluster/docker"

	"github.com/samalba/dockerclient"
)

var yaml2 = `

setup:
  image: plugins/drone-build
  commands:
    - go version
    - pwd
    - ls -la

clone:
  image: plugins/drone-git

build:
  image: golang:1.4
  command: [/drone/bin/build.sh]
  entrypoint: [/bin/bash]

compose:
  postgres:
    image: postgres:9.2

`

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

// IN: Config
// OUT: []Container

// IN: Req, Res
// OUT: error

// IN: Container, Engine

func do() {

	c := context{}
	c.request = &builder.Request{}
	c.response = &builder.Response{}
	c.engine = docker.New(nil)

	c.engine.Setup()
	defer c.engine.Teardown()

	err := build(c)
	if err != nil {
		c.response.ExitCode = 255
	}

	if c.response.ExitCode != 0 {
		err = deploy(c)
		if err != nil {
			c.response.ExitCode = 255
		}
	}

	err = notify(c)
	if err != nil {
		c.response.ExitCode = 255
	}
}

type context struct {
	request    *builder.Request
	response   *builder.Response
	engine     engine.Engine
	containers []*common.Container
}

func build(c context) error {
	b := builder.Builder{}
	defer b.Cancel()
	for _, container := range c.containers {
		switch container.Type {
		case "setup", "clone", "build":
			//b.Register(builder.Batch(c.environment, container))
		case "service":
			//b.Register(builder.Service(c.environment, container))
		}
	}
	return b.Build(c.response)
}

func deploy(c context) error {
	b := builder.Builder{}
	defer b.Cancel()
	for _, container := range c.containers {
		switch container.Type {
		case "publish", "deploy":
			//b.Register(builder.Batch(c.environment, container))
		}
	}
	return b.Build(c.response)
}

func notify(c context) error {
	b := builder.Builder{}
	defer b.Cancel()
	for _, container := range c.containers {
		switch container.Type {
		case "notify":
			//b.Register(builder.Batch(c.environment, container))
		}
	}
	return b.Build(c.response)
}

// ToContainers is a helper function that converts the build
// configuration to a list of containers.
func ToContainers() []*common.Container {
	config := struct {
		Setup *common.Image
		Clone *common.Image
		Build *common.Image

		Compose map[string]*common.Image
		Publish map[string]*common.Image
		Deploy  map[string]*common.Image
		Notify  map[string]*common.Image
	}{}

	//

	var containers []*common.Container

	for _, image := range config.Compose {
		containers = append(containers, toContainer(image, TypeService))
	}

	containers = append(containers, toSetup(config.Build))
	containers = append(containers, toClone(config.Clone))
	containers = append(containers, toBuild(config.Clone))

	for _, image := range config.Publish {
		containers = append(containers, toContainer(image, TypePublish))
	}
	for _, image := range config.Deploy {
		containers = append(containers, toContainer(image, TypeDeploy))
	}
	for _, image := range config.Notify {
		containers = append(containers, toContainer(image, TypeNotify))
	}

	return containers
}

func toClone(img *common.Image) *common.Container {
	return toContainer(&common.Image{
		Name:     "plugin/drone-git",
		UserData: img.UserData,
	}, TypeSetup)
}

func toSetup(i *common.Image) *common.Container {
	return toContainer(&common.Image{
		Name:     "plugin/drone-build",
		UserData: i.UserData,
	}, TypeSetup)
}

func toBuild(img *common.Image) *common.Container {
	img.UserData = nil
	img.Entrypoint = []string{"/bin/bash"}
	img.Cmd = []string{"/drone/bin/build.sh"}
	return toContainer(img, TypeSetup)
}

func toContainer(i *common.Image, t string) *common.Container {
	return &common.Container{
		Type:  t,
		Image: i,
	}
}

const (
	TypeSetup   = "setup"
	TypeClone   = "clone"
	TypeBuild   = "build"
	TypeService = "service"
	TypePublish = "publish"
	TypeDeploy  = "deploy"
	TypeNotify  = "notify"
)
