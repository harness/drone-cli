package main

import (
	"os"

	"github.com/drone/drone-cli/builder"
	"github.com/drone/drone-cli/builder/ambassador"
	"github.com/drone/drone-cli/common"
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
		client, err := ambassador.Create(&mockClient{})
		if err != nil {
			return
		}
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
