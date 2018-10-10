package exec

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/drone/envsubst"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-runtime/engine/docker"
	"github.com/drone/drone-runtime/runtime"
	"github.com/drone/drone-runtime/runtime/term"
	"github.com/drone/drone-yaml/yaml"
	"github.com/drone/drone-yaml/yaml/compiler"
	"github.com/drone/drone-yaml/yaml/compiler/transform"
	"github.com/drone/drone-yaml/yaml/converter"
	"github.com/drone/drone-yaml/yaml/linter"
	"github.com/drone/signal"

	"github.com/joho/godotenv"
	"github.com/mattn/go-isatty"
	"github.com/urfave/cli"
)

var tty = isatty.IsTerminal(os.Stdout.Fd())

// Command exports the exec command.
var Command = cli.Command{
	Name:      "exec",
	Usage:     "execute a local build",
	ArgsUsage: "[path/to/.drone.yml]",
	Action: func(c *cli.Context) {
		if err := exec(c); err != nil {
			log.Fatalln(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "pipeline",
			Usage: "Name of the pipelie to execute",
		},
		cli.StringSliceFlag{
			Name:  "include",
			Usage: "Name of steps to include",
		},
		cli.StringSliceFlag{
			Name:  "exclude",
			Usage: "Name of steps to exclude",
		},
		cli.StringFlag{
			Name:  "resume-at",
			Usage: "Name of start to resume at",
		},
		cli.BoolFlag{
			Name:  "clone",
			Usage: "enable the clone step",
		},
		cli.BoolFlag{
			Name:  "trusted",
			Usage: "build is trusted",
		},
		cli.DurationFlag{
			Name:  "timeout",
			Usage: "build timeout",
			Value: time.Hour,
		},
		cli.StringSliceFlag{
			Name:  "volume",
			Usage: "build volumes",
		},
		cli.StringSliceFlag{
			Name:  "network",
			Usage: "external networks",
		},
		cli.StringFlag{
			Name:  "secret-file",
			Usage: "secret file",
		},
		cli.StringSliceFlag{
			Name:  "privileged",
			Usage: "privileged plugins",
			Value: &cli.StringSlice{
				"plugins/docker",
				"plugins/gcr",
				"plugins/ecr",
			},
		},

		//
		// netrc parameters
		//
		cli.StringFlag{
			Name: "netrc-username",
		},
		cli.StringFlag{
			Name: "netrc-password",
		},
		cli.StringFlag{
			Name: "netrc-machine",
		},

		//
		// trigger parameters
		//

		cli.StringFlag{
			Name:  "branch",
			Usage: "branch name",
		},
		cli.StringFlag{
			Name:  "event",
			Usage: "build event name (push, pull_request, etc)",
		},
		cli.StringFlag{
			Name:  "instance",
			Usage: "instance hostname (e.g. drone.company.com)",
		},
		cli.StringFlag{
			Name:  "ref",
			Usage: "git reference",
		},
		cli.StringFlag{
			Name:  "repo",
			Usage: "git repository name (e.g. ocotcat/hello-world)",
		},
		cli.StringFlag{
			Name:  "deploy-to",
			Usage: "deployment target (e.g. production)",
		},
	},
}

func exec(c *cli.Context) error {
	file := c.Args().First()
	if file == "" {
		file = ".drone.yml"
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	environ := getEnv(c)
	dataS, err := envsubst.Eval(string(data), func(name string) string {
		return environ[name]
	})
	if err != nil {
		return err
	}

	// this code is temporarily in place to detect and convert
	// the legacy yaml configuration file to the new format.
	if converter.IsLegacy(dataS) {
		dataS, err = converter.ConvertString(dataS)
		if err != nil {
			return err
		}
	}

	manifest, err := yaml.ParseString(dataS)
	if err != nil {
		return err
	}

	var pipeline *yaml.Pipeline
	filter := c.String("pipeline")
	for _, resource := range manifest.Resources {
		v, ok := resource.(*yaml.Pipeline)
		if !ok {
			continue
		}
		if filter == "" || filter == v.Name {
			pipeline = v
			break
		}
	}
	if pipeline == nil {
		return errors.New("cannot find pipeline")
	}

	trusted := c.Bool("trusted")
	err = linter.Lint(pipeline, trusted)
	if err != nil {
		return err
	}

	// the user has the option to disable the git clone
	// if the pipeline is being executed on the local
	// codebase.
	if c.Bool("clone") == false {
		pipeline.Clone.Disable = true
	}

	comp := new(compiler.Compiler)
	comp.PrivilegedFunc = compiler.DindFunc(
		c.StringSlice("privileged"),
	)
	comp.SkipFunc = compiler.SkipFunc(
		compiler.SkipData{
			Branch:   environ["DRONE_BRANCH"],
			Event:    environ["DRONE_EVENT"],
			Instance: environ["DRONE_SYSTEM_HOST"],
			Ref:      environ["DRONE_COMMIT_REF"],
			Repo:     environ["DRONE_REPO"],
			Target:   environ["DRONE_DEPLOY_TO"],
		},
	)
	comp.TransformFunc = transform.Combine(
		transform.Include(
			c.StringSlice("include"),
		),
		transform.Exclude(
			c.StringSlice("exclude"),
		),
		transform.ResumeAt(
			c.String("resume-at"),
		),
		transform.WithAuths(
			toRegistry(
				c.StringSlice("registry"),
			),
		),
		transform.WithEnviron(
			readParams(
				c.String("env-file"),
			),
		),
		transform.WithEnviron(environ),
		transform.WithLables(nil),
		transform.WithLimits(0, 0),
		transform.WithNetrc(
			c.String("netrc-machine"),
			c.String("netrc-username"),
			c.String("netrc-password"),
		),
		transform.WithNetworks(
			c.StringSlice("network"),
		),
		transform.WithProxy(),
		transform.WithSecrets(
			readParams(
				c.String("env-file"),
			),
		),
		transform.WithVolumes(
			toVolumes(
				c.StringSlice("volume"),
			),
		),
	)
	ir := comp.Compile(pipeline)

	// the user has the option to disable the git clone
	// if the pipeline is being executed on the local
	// codebase.
	if c.Bool("clone") == false {
		pwd, _ := os.Getwd()
		mountWorkspace(ir, pwd)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		c.Duration("timeout"),
	)
	ctx = signal.WithContext(ctx)
	defer cancel()

	// creates a docker-based engine. eventually we will
	// include the kubernetes and vmware fusion engines.
	engine, err := docker.NewEnv()
	if err != nil {
		return err
	}

	// creates a hook to print the step output to stdout,
	// with per-step color coding if a tty.
	hooks := &runtime.Hook{}
	hooks.GotLine = term.WriteLine(os.Stdout)
	if tty {
		hooks.GotLine = term.WriteLinePretty(os.Stdout)
	}

	return runtime.New(
		runtime.WithEngine(engine),
		runtime.WithConfig(ir),
		runtime.WithHooks(hooks),
	).Run(ctx)
}

// helper function converts a slice of colon-separated
// volumes to a map.
func toVolumes(items []string) map[string]string {
	set := map[string]string{}
	for _, item := range items {
		parts := strings.Split(item, ":")
		if len(parts) != 2 {
			key := parts[0]
			val := parts[1]
			set[key] = val
		}
	}
	return set
}

// helper function converts a slice of urls to a slice
// of docker registry credentials.
func toRegistry(items []string) []*engine.DockerAuth {
	auths := []*engine.DockerAuth{}
	for _, item := range items {
		uri, err := url.Parse(item)
		if err != nil {
			continue // skip invalid
		}
		host := uri.Host
		user := uri.User.Username()
		pass, _ := uri.User.Password()
		auths = append(auths, &engine.DockerAuth{
			Address:  host,
			Username: user,
			Password: pass,
		})
	}
	return auths
}

// helper function reads secrets from a key-value file.
func readParams(path string) map[string]string {
	data, _ := godotenv.Read(path)
	return data
}
