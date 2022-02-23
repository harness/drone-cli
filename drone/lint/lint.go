package lint

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/drone-runners/drone-runner-docker/engine/linter"
	"github.com/drone-runners/drone-runner-docker/engine/resource"
	"github.com/drone/drone-go/drone"
	"github.com/drone/envsubst"
	"github.com/drone/runner-go/manifest"
	"github.com/urfave/cli"
)

// Command exports the linter command.
var Command = cli.Command{
	Name:      "lint",
	Usage:     "lint the yaml file, checks for yaml errors",
	ArgsUsage: "<source>",
	Action:    lint,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "trusted",
			Usage: "is the yaml trustable",
		},
	},
}

type Flags struct {
	Build  drone.Build
	Netrc  drone.Netrc
	Repo   drone.Repo
	Stage  drone.Stage
	System drone.System
}

func lint(c *cli.Context) error {
	f := new(Flags)
	f.Repo.Trusted = c.Bool("trusted")
	var envs map[string]string

	path := c.Args().First()
	if path == "" {
		path = ".drone.yml"
	}

	rawSource, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	// string substitution function ensures that string replacement variables are escaped and quoted if they contain newlines.
	subf := func(k string) string {
		v := envs[k]
		if strings.Contains(v, "\n") {
			v = fmt.Sprintf("%q", v)
		}
		return v
	}
	// evaluates string replacement expressions and returns an update configuration.
	config, err := envsubst.Eval(string(rawSource), subf)
	if err != nil {
		return err
	}
	// parse into manifests
	inputManifests, err := manifest.ParseString(config)
	if err != nil {
		return err
	}
	for _, iter := range inputManifests.Resources {
		if iter.GetType() == "docker" {
			resource, err := resource.Lookup(iter.GetName(), inputManifests)
			if err != nil {
				return err
			}
			// lint the resource and return an error if any linting rules are broken
			lint := linter.New()
			err = lint.Lint(resource, &f.Repo)
			if err != nil {
				return err
			}
			fmt.Printf("%v\n", iter)
		}
	}
	// now we can check the pipeline dependencies
	// get a list of all the pipelines
	allStages := map[string]struct{}{}
	for _, iter := range inputManifests.Resources {
		allStages[iter.GetName()] = struct{}{}
	}
	// we need to parse the file again into raw resources to access the dependencies
	inputRawResources, err := manifest.ParseRawFile(path)
	if err != nil {
		return err
	}
	for _, iter := range inputRawResources {
		for _, dep := range iter.Deps {
			if _, ok := allStages[dep]; !ok {
				return fmt.Errorf("Pipeline stage '%s' declares invalid dependency '%s'", iter.Name, dep)
			}
		}
	}
	return nil
}
