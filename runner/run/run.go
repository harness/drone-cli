package main

import (
	"fmt"
	"os"

	"github.com/drone/drone-cli/builder"
	"github.com/drone/drone-cli/builder/ambassador"
	"github.com/drone/drone-cli/builder/runner"
	"github.com/drone/drone-cli/common"
	"github.com/drone/drone-cli/common/config"
)

func main() {

	matrix, err := config.ParseMatrix(testYaml)
	if err != nil {
		println(err.Error())
		return
	}

	client, err := newMockClient() //dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, conf := range matrix {
		config.Transform(conf)
		build.Client = client
		build.Config = conf
		run(build)
	}
}

func run(build *builder.Build) {

	amb, err := ambassador.Create(build.Client)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer amb.Destroy()
	build.Client = amb // TODO remove this

	// response writer
	res := builder.NewResultWriter(os.Stdout)

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
	n := runner.Notifier(build)
	defer n.Cancel()
	err = n.Build(res)
	if err != nil {
		fmt.Println(err)
	}
}

var build = &builder.Build{
	Repo: &common.Repo{
		Remote: "github.com",
		Host:   "github.com",
		Owner:  "bradrydzewski",
	},
	Clone: &common.Clone{
		Dir:    "/drone/src/github.com/drone/drone",
		Sha:    "4fbcc1dd41c5e2792c034d31e350f521890ad723",
		Branch: "master",
		Remote: "git://github.com/drone/drone.git",
	},
}

var testYaml = `
build:
  image: golang:1.4.2
  commands:
    - ls -la /drone/src/github.com/drone/drone
    - go version

notify:
  slack:
    channel: dev
    usernae: drone

compose:
  mysql:
    image: bradrydzewski/mysql:5.5
  postgres:
    image: bradrydzewski/postgres:9.1
`

// compose:
//   database:
// 	  image: postgres:9.2

// matrix:
//   go_version:
//     - 1.3.3
//     - 1.4.2

var droneYaml = `
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
