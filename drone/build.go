package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/drone/drone-cli/builder"
	"github.com/drone/drone-cli/builder/docker"
	"github.com/drone/drone-cli/common"
	"github.com/drone/drone-cli/parser"
	"github.com/drone/drone-cli/parser/inject"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/samalba/dockerclient"
)

const EXIT_STATUS = 1

// NewBuildCommand returns the CLI command for "build".
func NewBuildCommand() cli.Command {
	return cli.Command{
		Name:  "build",
		Usage: "run a local build",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "i",
				Value: "",
				Usage: "identify file injected in the container",
			},
			cli.BoolFlag{
				Name:  "p",
				Usage: "runs drone build in a privileged container",
			},
			cli.BoolFlag{
				Name:  "deploy",
				Usage: "runs drone build with deployments enabled",
			},
			cli.BoolFlag{
				Name:  "publish",
				Usage: "runs drone build with publishing enabled",
			},
			cli.StringFlag{
				Name:   "docker-host",
				Value:  "unix:///var/run/docker.sock",
				Usage:  "docker daemon address",
				EnvVar: "DOCKER_HOST",
			},
			cli.StringFlag{
				Name:  "docker-cert",
				Value: getCert(),
				Usage: "docker daemon tls certificate",
			},
			cli.StringFlag{
				Name:  "docker-key",
				Value: getKey(),
				Usage: "docker daemon tls key",
			},
		},
		Action: func(c *cli.Context) {
			buildCommandFunc(c)
		},
	}
}

// buildCommandFunc executes the "build" command.
func buildCommandFunc(c *cli.Context) {
	var privileged = c.Bool("p")
	var identity = c.String("i")
	var deploy = c.Bool("deploy")
	var publish = c.Bool("publish")
	var path string

	var dockerhost = c.String("docker-host")
	var dockercert = c.String("docker-cert")
	var dockerkey = c.String("docker-key")

	// the path is provided as an optional argument that
	// will otherwise default to $PWD/.drone.yml
	if len(c.Args()) > 0 {
		path = c.Args()[0]
	}

	switch len(path) {
	case 0:
		path, _ = os.Getwd()
		path = filepath.Join(path, ".drone.yml")
	default:
		path = filepath.Clean(path)
		path, _ = filepath.Abs(path)
		path = filepath.Join(path, ".drone.yml")
	}

	var exit = run(path, identity, dockerhost, dockercert, dockerkey, publish, deploy, privileged)
	os.Exit(exit)
}

// TODO this has gotten a bit out of hand. refactor input params
func run(path, identity, dockerhost, dockercert, dockerkey string, publish, deploy, privileged bool) int {
	// dockerClient, err := docker.NewHostCertFile(dockerhost, dockercert, dockerkey)
	// if err != nil {
	// 	log.Err(err.Error())
	// 	return EXIT_STATUS, err
	// }

	// parse the private environment variables
	envs := getParamMap("DRONE_ENV_")

	// parse the drone.yml file
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return EXIT_STATUS
	}
	yml := inject.Inject(string(raw), envs)

	matrix, err := parser.Parse(yml)
	if err != nil {
		return EXIT_STATUS
	}

	// get the repository root directory
	parent_dir := filepath.Dir(path)
	dir := filepath.Dir(path)

	// does the local repository match the
	// $GOPATH/src/{package} pattern? This is
	// important so we know the target location
	// where the code should be copied inside
	// the container.
	if gopath, ok := getRepoPath(dir); ok {
		dir = gopath

	} else if gopath, ok := getGoPath(dir); ok {
		// in this case we found a GOPATH and
		// reverse engineered the package path
		dir = gopath

	} else {
		// otherwise just use directory name
		dir = filepath.Base(dir)
	}

	// this is where the code gets uploaded to the container
	// TODO move this code to the build package
	dir = filepath.Join("/drone/src", filepath.Clean(dir))

	// ssh key to import into container
	var key []byte
	if len(identity) != 0 {
		key, err = ioutil.ReadFile(identity)
		if err != nil {
			fmt.Printf("[Error] Could not find or read identity file %s\n", identity)
			return EXIT_STATUS
		}
	}

	//
	//
	//

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

		// /home/brad/gocode/src/github.com/garyburd/redigo:/drone/src/github.com/garyburd/redigo

		conf.Build.Volumes = append(conf.Setup.Volumes, parent_dir+":"+dir)
		conf.Clone = nil

		//client := &mockClient{}
		client, _ := dockerclient.NewDockerClient(dockerhost, nil)
		ambassador, err := docker.NewAmbassador(client)
		if err != nil {
			return EXIT_STATUS
		}

		c := Context{}
		c.builder = builder.Load(conf)
		c.build = builder.NewB(ambassador, os.Stdout)
		c.build.Repo = &common.Repo{}
		c.build.Commit = &common.Commit{}
		c.build.Clone = &common.Clone{Dir: dir, Keypair: &common.Keypair{Private: string(key)}}
		c.config = conf
		c.client = ambassador

		contexts = append(contexts, &c)
	}

	// run the builds
	var exit int
	for _, c := range contexts {
		log.Printf("starting build %s", c.config.Axis)
		err := c.builder.RunBuild(c.build)
		if err != nil {
			c.build.Exit(255)
			// TODO need a 255 exit code if the build errors
		}
		if c.build.ExitCode() != 0 {
			exit = c.build.ExitCode()
		}
	}

	// run the deploy steps
	if exit == 0 {
		for _, c := range contexts {
			if !c.builder.HasDeploy() {
				continue
			}
			log.Printf("starting post-build tasks %s", c.config.Axis)
			err := c.builder.RunDeploy(c.build)
			if err != nil {
				c.build.Exit(255)
				// TODO need a 255 exit code if the build errors
			}
			if c.build.ExitCode() != 0 {
				exit = c.build.ExitCode()
			}
		}
	}

	// run the notify steps
	for _, c := range contexts {
		if !c.builder.HasNotify() {
			continue
		}
		log.Printf("staring notification tasks %s", c.config.Axis)
		c.builder.RunNotify(c.build)
		break
	}

	log.Println("build complete")
	for _, c := range contexts {
		log.WithField("exit_code", c.build.ExitCode()).Infoln(c.config.Axis)
	}

	return exit
}

func getCert() string {
	if os.Getenv("DOCKER_CERT_PATH") != "" && os.Getenv("DOCKER_TLS_VERIFY") == "1" {
		return filepath.Join(os.Getenv("DOCKER_CERT_PATH"), "cert.pem")
	} else {
		return ""
	}
}

func getKey() string {
	if os.Getenv("DOCKER_CERT_PATH") != "" && os.Getenv("DOCKER_TLS_VERIFY") == "1" {
		return filepath.Join(os.Getenv("DOCKER_CERT_PATH"), "key.pem")
	} else {
		return ""
	}
}

func init() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&formatter{})
}

type Context struct {
	build   *builder.B
	builder *builder.Builder
	config  *common.Config

	client *docker.Ambassador
}

type formatter struct {
	nocolor bool
}

func (f *formatter) Format(entry *log.Entry) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("\033[2m")
	buf.WriteString("[drone]")

	for k, v := range entry.Data {
		if k != "exit_code" {
			continue
		}

		if v == 0 {
			buf.WriteString("\033[1;32m SUCCESS\033[0m")
		} else {
			buf.WriteString("\033[1;31m FAILURE\033[0m")
		}
	}

	buf.WriteByte(' ')
	buf.WriteString(entry.Message)
	buf.WriteByte(' ')

	for k, v := range entry.Data {
		buf.WriteString(
			fmt.Sprintf("%s=%v", k, v),
		)
	}

	buf.WriteString("\033[0m")
	buf.WriteByte('\n')
	return buf.Bytes(), nil
}
