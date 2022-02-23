package exec2

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/drone-runners/drone-runner-docker/engine"
	"github.com/drone-runners/drone-runner-docker/engine/compiler"
	"github.com/drone-runners/drone-runner-docker/engine/linter"
	"github.com/drone-runners/drone-runner-docker/engine/resource"

	"github.com/drone/drone-go/drone"
	"github.com/drone/envsubst"
	"github.com/drone/runner-go/environ"
	"github.com/drone/runner-go/environ/provider"
	"github.com/drone/runner-go/logger"
	"github.com/drone/runner-go/manifest"
	"github.com/drone/runner-go/pipeline"
	"github.com/drone/runner-go/pipeline/runtime"
	"github.com/drone/runner-go/pipeline/streamer/console"
	"github.com/drone/runner-go/registry"
	"github.com/drone/runner-go/secret"
	"github.com/drone/signal"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var nocontext = context.Background()

// Command exports the exec command.
var Command = cli.Command{
	Name:      "exec2",
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
			Usage: "Name of the pipeline to execute",
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
			Name:  "registry",
			Usage: "registry file",
		},
		cli.StringFlag{
			Name:  "secret-file",
			Usage: "secret file, define values that can be used with from_secret",
		},
		cli.StringFlag{
			Name:  "env-file",
			Usage: "env file",
		},
		cli.StringSliceFlag{
			Name:  "privileged",
			Usage: "privileged plugins",
			Value: &cli.StringSlice{
				"plugins/docker",
				"plugins/acr",
				"plugins/ecr",
				"plugins/gcr",
				"plugins/heroku",
			},
		},
		// netrc parameters
		cli.StringFlag{
			Name: "netrc-username",
		},
		cli.StringFlag{
			Name: "netrc-password",
		},
		cli.StringFlag{
			Name: "netrc-machine",
		},
		// trigger parameters
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
		// cli.StringFlag{
		// 	Name:  "ref",
		// 	Usage: "git reference",
		// }, NOT NEEDED
		// cli.StringFlag{
		// 	Name:  "sha",
		// 	Usage: "git sha",
		// },
		cli.StringFlag{
			Name:  "repo",
			Usage: "git repository name (e.g. octocat/hello-world)",
		},
		cli.StringFlag{
			Name:  "deploy-to",
			Usage: "deployment target (e.g. production)",
		},
	},
}

func exec(cliContext *cli.Context) error {
	// lets do our mapping from CLI flags to an execCommand struct
	commy := mapOldToExecCommand(cliContext)

	rawsource, err := ioutil.ReadFile(commy.Source)
	if err != nil {
		return err
	}
	envs := environ.Combine(
		getEnv(cliContext),
		environ.System(commy.System),
		environ.Repo(commy.Repo),
		environ.Build(commy.Build),
		environ.Stage(commy.Stage),
		environ.Link(commy.Repo, commy.Build, commy.System),
		commy.Build.Params,
	)

	// string substitution function ensures that string
	// replacement variables are escaped and quoted if they
	// contain newlines.
	subf := func(k string) string {
		v := envs[k]
		if strings.Contains(v, "\n") {
			v = fmt.Sprintf("%q", v)
		}
		return v
	}

	// evaluates string replacement expressions and returns an
	// update configuration.
	config, err := envsubst.Eval(string(rawsource), subf)
	if err != nil {
		return err
	}

	// parse and lint the configuration.
	manifest, err := manifest.ParseString(config)
	if err != nil {
		return err
	}

	// a configuration can contain multiple pipelines.
	// get a specific pipeline resource for execution.
	if commy.Stage.Name == "" {
		fmt.Println("No stage specified, assuming 'default'")
	}

	res, err := resource.Lookup(commy.Stage.Name, manifest)
	if err != nil {
		return fmt.Errorf("Stage '%s' not found in build file : %s", commy.Stage.Name, err)
	}

	// lint the pipeline and return an error if any
	// linting rules are broken
	lint := linter.New()
	err = lint.Lint(res, commy.Repo)
	if err != nil {
		return err
	}

	// compile the pipeline to an intermediate representation.
	comp := &compiler.Compiler{
		Environ:    provider.Static(commy.Environ),
		Labels:     commy.Labels,
		Resources:  commy.Resources,
		Tmate:      commy.Tmate,
		Privileged: append(commy.Privileged, compiler.Privileged...),
		Networks:   commy.Networks,
		Volumes:    commy.Volumes,
		Secret:     secret.StaticVars(commy.Secrets),
		Registry: registry.Combine(
			registry.File(commy.Config),
		),
	}

	// when running a build locally cloning is always
	// disabled in favor of mounting the source code
	// from the current working directory.
	if !commy.Clone {
		comp.Mount, _ = os.Getwd()
	}

	args := runtime.CompilerArgs{
		Pipeline: res,
		Manifest: manifest,
		Build:    commy.Build,
		Netrc:    commy.Netrc,
		Repo:     commy.Repo,
		Stage:    commy.Stage,
		System:   commy.System,
	}
	spec := comp.Compile(nocontext, args).(*engine.Spec)

	// include only steps that are in the include list,
	// if the list in non-empty.
	if len(commy.Include) > 0 {
	I:
		for _, step := range spec.Steps {
			if step.Name == "clone" {
				continue
			}
			for _, name := range commy.Include {
				if step.Name == name {
					continue I
				}
			}
			step.RunPolicy = runtime.RunNever
		}
	}
	// exclude steps that are in the exclude list, if the list in non-empty.
	if len(commy.Exclude) > 0 {
	E:
		for _, step := range spec.Steps {
			if step.Name == "clone" {
				continue
			}
			for _, name := range commy.Exclude {
				if step.Name == name {
					step.RunPolicy = runtime.RunNever
					continue E
				}
			}
		}
	}
	// resume at a specific step
	if cliContext.String("resume-at") != "" {
		for _, step := range spec.Steps {
			if step.Name == cliContext.String("resume-at") {
				break
			}
			if step.Name == "clone" {
				continue
			}
			for _, name := range commy.Exclude {
				if step.Name == name {
					step.RunPolicy = runtime.RunNever
					continue
				}
			}
		}
	}
	// create a step object for each pipeline step.
	for _, step := range spec.Steps {
		if step.RunPolicy == runtime.RunNever {
			continue
		}
		commy.Stage.Steps = append(commy.Stage.Steps, &drone.Step{
			StageID:   commy.Stage.ID,
			Number:    len(commy.Stage.Steps) + 1,
			Name:      step.Name,
			Status:    drone.StatusPending,
			ErrIgnore: step.ErrPolicy == runtime.ErrIgnore,
		})
	}

	// configures the pipeline timeout.
	timeout := time.Duration(commy.Repo.Timeout) * time.Minute
	ctx, cancel := context.WithTimeout(nocontext, timeout)
	defer cancel()

	// listen for operating system signals and cancel execution when received.
	ctx = signal.WithContextFunc(ctx, func() {
		println("received signal, terminating process")
		cancel()
	})

	state := &pipeline.State{
		Build:  commy.Build,
		Stage:  commy.Stage,
		Repo:   commy.Repo,
		System: commy.System,
	}

	// enable debug logging
	logrus.SetLevel(logrus.WarnLevel)
	if commy.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if commy.Trace {
		logrus.SetLevel(logrus.TraceLevel)
	}
	logger.Default = logger.Logrus(
		logrus.NewEntry(
			logrus.StandardLogger(),
		),
	)

	engine, err := engine.NewEnv(engine.Opts{})
	if err != nil {
		return err
	}

	err = runtime.NewExecer(
		pipeline.NopReporter(),
		console.New(commy.Pretty),
		pipeline.NopUploader(),
		engine,
		commy.Procs,
	).Exec(ctx, spec, state)

	if err != nil {
		dump(state)
		return err
	}
	switch state.Stage.Status {
	case drone.StatusError, drone.StatusFailing, drone.StatusKilled:
		os.Exit(1)
	}
	return nil
}

func dump(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}
