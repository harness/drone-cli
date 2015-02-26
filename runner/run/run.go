package main

import (
	"os"

	"github.com/drone/drone-cli/config"
	"github.com/drone/drone-cli/runner"
	"github.com/samalba/dockerclient"
)

func main() {

	conf := config.Config{}
	conf.Image = "library/golang:1.4.1-cross"
	conf.Script = []string{
		"ls /drone/src/github.com/drone/drone",
		// "add-apt-repository ppa:git-core/ppa 1> /dev/null 2> /dev/null",
		// "sudo apt-get update 1> /dev/null 2> /dev/null",
		// "sudo apt-get update 1> /dev/null 2> /dev/null",
		// "sudo apt-get -y install git zip libsqlite3-dev sqlite3 rpm 1> /dev/null 2> /dev/null",
		// "make docker",
		"make deps",
		// "make test",
		// "make test_postgres",
		// "make test_mysql",
	}
	conf.Services = []string{
		"bradrydzewski/mysql:5.5",
		"bradrydzewski/postgres:9.1",
	}

	clone := runner.Clone{}
	clone.Dir = "/drone/src/github.com/drone/drone"
	clone.Branch = "master"
	clone.Sha = "4fbcc1dd41c5e2792c034d31e350f521890ad723"
	clone.Remote = "git://github.com/drone/drone.git"

	var err error
	req := runner.Request{}
	req.Clone = &clone
	req.Config = &conf
	req.Client, err = dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
	if err != nil {
		println(err.Error())
		return
	}

	resp := runner.Response{}
	resp.Writer = os.Stdout
	runner.Run(&req, &resp)
	os.Exit(resp.ExitCode)
}
