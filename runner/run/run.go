package main

import (
	"os"

	"github.com/drone/drone-cli/builder"
	"github.com/drone/drone-cli/builder/docker"
	"github.com/drone/drone-cli/common"
	"github.com/drone/drone-cli/parser"

	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&formatter{})
}

type Context struct {
	build   *builder.B
	builder *builder.Builder
	config  *common.Config

	client *docker.Ambassador
}

func main() {
	// parse the build matrix
	matrix, err := parser.Parse(testYaml)
	if err != nil {
		println(err.Error())
		return
	}

	var contexts []*Context

	// must cleanup after our build
	defer func() {
		for _, c := range contexts {
			c.build.RemoveAll()
			c.client.Destroy()
		}
	}()

	// list of builds and builders for each item
	// in the matrix
	for _, conf := range matrix {
		//client := &mockClient{}
		client, _ := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
		ambassador, err := docker.NewAmbassador(client)
		if err != nil {
			return
		}

		c := Context{}
		c.builder = builder.Load(conf)
		c.build = builder.NewB(ambassador, os.Stdout)
		c.build.Repo = repo
		c.build.Clone = clone
		c.build.Commit = commit
		c.config = conf
		c.client = ambassador

		contexts = append(contexts, &c)
	}

	// run the builds
	for _, c := range contexts {
		log.Printf("starting build %s", c.config.Axis)
		err := c.builder.RunBuild(c.build)
		if err != nil {
			// TODO need a 255 exit code if the build errors
		}
	}

	// run the deploy steps
	// run the notify steps

	log.Println("")
	for _, c := range contexts {
		log.Printf(" ✓ %s", c.config.Axis)
	}
	log.Println("")
}

var repo = &common.Repo{
	Remote: "github.com",
	Host:   "github.com",
	Owner:  "bradrydzewski",
}
var clone = &common.Clone{
	Dir:    "/drone/src/github.com/garyburd/redigo",
	Sha:    "535138d7bcd717d6531c701ef5933d98b1866257",
	Branch: "master",
	Remote: "git://github.com/garyburd/redigo.git",
}

var commit = &common.Commit{}

var testYaml = `
build:
  image: golang:$$go_version
  environment:
    - GOPATH=/drone
  commands:
    - cd redis
    - go version
    - go build
    - go test -v

compose:
  redis:
    image: redis

matrix:
  go_version:
    - 1.3.3
    - 1.4.2
`

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
