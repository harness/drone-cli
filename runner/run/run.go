package main

import (
	"fmt"
	"os"

	"github.com/drone/drone-cli/builder"
	"github.com/drone/drone-cli/builder/ambassador"
	"github.com/drone/drone-cli/builder/runner"
	"github.com/drone/drone-cli/common"
	"github.com/drone/drone-cli/common/config"
	"github.com/samalba/dockerclient"
)

var yaml = `
build:
  image: golang:$$go_version
  script:
    - add-apt-repository ppa:git-core/ppa 1> /dev/null 2> /dev/null
    - apt-get update 1> /dev/null 2> /dev/null
    - apt-get update 1> /dev/null 2> /dev/null
    - apt-get -y install git zip libsqlite3-dev sqlite3 rpm 1> /dev/null 2> /dev/null
    - make docker
    - make deps

compose:
  mysql:
    image: bradrydzewski/mysql:5.5
  postgres:
    image: bradrydzewski/postgres:9.1

matrix:
  go_version:
    - 1.3.3
    - 1.4.2
`

var yamlAlt = `
clone:
  image: plugin/drone-git

build:
  image: golang:1.4.2
  commands:
    - ls -la /drone/src/github.com/drone/drone
    - go version

compose:
  database:
	  image: postgres:9.2

#matrix:
#  go_version:
#    - 1.3.3
#    - 1.4.2
`

func main() {

	matrix, err := config.ParseMatrix(yamlAlt)
	if err != nil {
		println(err.Error())
		return
	}

	clone := &common.Clone{}
	clone.Dir = "/drone/src/github.com/drone/drone"
	clone.Branch = "master"
	clone.Sha = "4fbcc1dd41c5e2792c034d31e350f521890ad723"
	clone.Remote = "git://github.com/drone/drone.git"

	// client
	client, err := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// ambassador
	amb, err := ambassador.Create(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	// response writer
	res := builder.NewResultWriter(os.Stdout)

	// context
	build := &builder.Build{
		Repo:   nil,
		Clone:  clone,
		Commit: nil,
		Config: matrix[0],
		Client: amb,
	}

	// builder
	b := runner.Builder(build)
	defer b.Cancel()
	err = b.Build(res)
	if err != nil {
		fmt.Println(err)
	}

	// deployer
	d := runner.Deployer(build)
	defer d.Cancel()
	err = d.Build(res)
	if err != nil {
		fmt.Println(err)
	}

	// notifier
	n := runner.Deployer(build)
	defer n.Cancel()
	err = n.Build(res)
	if err != nil {
		fmt.Println(err)
	}

	// for _, conf := range matrix {
	// 	req := runner.Request{}
	// 	req.Clone = clone
	// 	req.Config = conf
	// 	req.Client, err = dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
	// 	if err != nil {
	// 		println(err.Error())
	// 		return
	// 	}
	//
	// 	resp := runner.Response{}
	// 	resp.Writer = os.Stdout
	// 	runner.Run(&req, &resp)
	//
	// 	if resp.ExitCode != 0 {
	// 		os.Exit(resp.ExitCode)
	// 	}
	//
	// }
}
