package exec

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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

		if v.Type != "" && v.Type != "docker" {
			return fmt.Errorf("pipeline type (%s) is not supported with 'drone exec'", v.Type)
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
	transforms := []func(*engine.Spec){
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
				c.String("secret-file"),
			),
		),
		transform.WithVolumeSlice(
			c.StringSlice("volume"),
		),
	}
	var pipelineFQN, pwd string
	if c.Bool("clone") == false {
		pwd, _ = os.Getwd()
		comp.WorkspaceMountFunc = compiler.MountHostWorkspace
		comp.WorkspaceFunc = compiler.CreateHostWorkspace(pwd)
		//a unique name pattern for pipeline with pipeline dir and name
		pipelineFQN = fmt.Sprintf("%s~~%s", strings.ReplaceAll(pwd, "/", "-"), pipeline.Name)
	}
	comp.TransformFunc = transform.Combine(transforms...)
	ir := comp.Compile(pipeline)

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
	hooks.BeforeEach = func(s *runtime.State) error {
		//add this label only for Drone exec usage
		if pwd != "" {
			s.Step.Metadata.Labels["io.drone.pipeline.dir"] = pwd
			s.Step.Metadata.Labels["io.drone.pipeline.name"] = pipelineFQN
		}
		s.Step.Envs["CI_BUILD_STATUS"] = "success"
		s.Step.Envs["CI_BUILD_STARTED"] = strconv.FormatInt(s.Runtime.Time, 10)
		s.Step.Envs["CI_BUILD_FINISHED"] = strconv.FormatInt(time.Now().Unix(), 10)
		s.Step.Envs["DRONE_BUILD_STATUS"] = "success"
		s.Step.Envs["DRONE_BUILD_STARTED"] = strconv.FormatInt(s.Runtime.Time, 10)
		s.Step.Envs["DRONE_BUILD_FINISHED"] = strconv.FormatInt(time.Now().Unix(), 10)

		s.Step.Envs["CI_JOB_STATUS"] = "success"
		s.Step.Envs["CI_JOB_STARTED"] = strconv.FormatInt(s.Runtime.Time, 10)
		s.Step.Envs["CI_JOB_FINISHED"] = strconv.FormatInt(time.Now().Unix(), 10)
		s.Step.Envs["DRONE_JOB_STATUS"] = "success"
		s.Step.Envs["DRONE_JOB_STARTED"] = strconv.FormatInt(s.Runtime.Time, 10)
		s.Step.Envs["DRONE_JOB_FINISHED"] = strconv.FormatInt(time.Now().Unix(), 10)

		if s.Runtime.Error != nil {
			s.Step.Envs["CI_BUILD_STATUS"] = "failure"
			s.Step.Envs["CI_JOB_STATUS"] = "failure"
			s.Step.Envs["DRONE_BUILD_STATUS"] = "failure"
			s.Step.Envs["DRONE_JOB_STATUS"] = "failure"
		}
		return nil
	}

	hooks.GotLine = term.WriteLine(os.Stdout)
	if tty {
		hooks.GotLine = term.WriteLinePretty(
			colorable.NewColorableStdout(),
		)
	}

	return runtime.New(
		runtime.WithEngine(engine),
		runtime.WithConfig(ir),
		runtime.WithHooks(hooks),
	).Run(ctx)
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
		user := uri.User.Username()
		pass, _ := uri.User.Password()
		uri.User = nil
		auths = append(auths, &engine.DockerAuth{
			Address:  uri.String(),
			Username: user,
			Password: pass,
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
