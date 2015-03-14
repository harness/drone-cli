package main

import (
	"os"

	"github.com/drone/drone-cli/common"
	builder "github.com/drone/drone-cli/compiler"
	"github.com/drone/drone-cli/parser"

	log "github.com/Sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&formatter{})
}

func main() {
	matrix, err := parser.Parse(testYaml)
	if err != nil {
		println(err.Error())
		return
	}

	client, err := newMockClient() //dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
	if err != nil {
		log.Errorln(err)
		return
	}
	// TODO DEFER AMABASSADOR DESTROY

	var builds []*builder.B
	var builders []*builder.Builder

	// must cleanup after our build
	defer func() {
		for _, build := range builds {
			build.RemoveAll()
		}
	}()

	// list of builds and builders for each item
	// in the matrix
	for _, conf := range matrix {
		b := builder.NewB(client, os.Stdout)
		b.Repo = repo
		b.Clone = clone
		b.Config = conf
		builds = append(builds, b)
		builders = append(builders, builder.Load(conf))
	}

	// run the builds
	for i, builder := range builders {
		log.Printf("starting build %s", builds[i].Config.Axis)
		err := builder.RunBuild(builds[i])
		if err != nil {
			// TODO need a 255 exit code if the build errors
		}
	}

	// run the deploy steps
	// run the notify steps

	log.Println("")
	for _, b := range builds {
		log.Printf(" ✓ %s", b.Config.Axis)
	}
	log.Println("")
}

//
// func run(build *builder.Build) {
//
// 	amb, err := ambassador.Create(build.Client)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer amb.Destroy()
// 	build.Client = amb // TODO remove this
//
// 	// response writer
// 	res := builder.NewResultWriter(os.Stdout)
//
// 	// builder
// 	b := runner.Builder(build)
// 	defer b.Cancel()
// 	err = b.Build(res)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
//
// 	// deployer
// 	d := runner.Deployer(build)
// 	defer d.Cancel()
// 	err = d.Build(res)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
//
// 	// notifier
// 	n := runner.Notifier(build)
// 	defer n.Cancel()
// 	err = n.Build(res)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

var repo = &common.Repo{
	Remote: "github.com",
	Host:   "github.com",
	Owner:  "bradrydzewski",
}
var clone = &common.Clone{
	Dir:    "/drone/src/github.com/drone/drone",
	Sha:    "4fbcc1dd41c5e2792c034d31e350f521890ad723",
	Branch: "master",
	Remote: "git://github.com/drone/drone.git",
}

var testYaml = `
build:
  image: golang:$$go_version
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

matrix:
  go_version:
    - 1.3.3
    - 1.4.2
`

// compose:
//   database:
// 	  image: postgres:9.2

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
