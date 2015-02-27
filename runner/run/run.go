package main

import (
	"os"

	"github.com/drone/drone-cli/common"
	"github.com/drone/drone-cli/common/config"
	"github.com/drone/drone-cli/runner"
	"github.com/samalba/dockerclient"
)

var yaml = `
image: golang:$$go_version
script:
  - add-apt-repository ppa:git-core/ppa 1> /dev/null 2> /dev/null
  - apt-get update 1> /dev/null 2> /dev/null
  - apt-get update 1> /dev/null 2> /dev/null
  - apt-get -y install git zip libsqlite3-dev sqlite3 rpm 1> /dev/null 2> /dev/null
  - make docker
  - make deps

services:
  - bradrydzewski/mysql:5.5
  - bradrydzewski/postgres:9.1

matrix:
  go_version:
    - 1.3.3
    - 1.4.2
`

var yaml_alt = `
image: golang:$$go_version
script:
  - ls -la /drone/src/github.com/drone/drone
  - go version

matrix:
  go_version:
    - 1.3.3
    - 1.4.2
`

func main() {

	matrix, err := config.ParseMatrix(yaml_alt)
	if err != nil {
		println(err.Error())
		return
	}

	clone := &common.Clone{}
	clone.Dir = "/drone/src/github.com/drone/drone"
	clone.Branch = "master"
	clone.Sha = "4fbcc1dd41c5e2792c034d31e350f521890ad723"
	clone.Remote = "git://github.com/drone/drone.git"

	for _, conf := range matrix {
		req := runner.Request{}
		req.Clone = clone
		req.Config = conf
		req.Client, err = dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
		if err != nil {
			println(err.Error())
			return
		}

		resp := runner.Response{}
		resp.Writer = os.Stdout
		runner.Run(&req, &resp)

		if resp.ExitCode != 0 {
			os.Exit(resp.ExitCode)
		}

	}
}
