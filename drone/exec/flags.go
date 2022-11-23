// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package exec

import (
	"strings"

	"github.com/drone-runners/drone-runner-docker/engine/compiler"
	"github.com/drone/drone-go/drone"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

// Flags maps
type Flags struct {
	Build  *drone.Build
	Netrc  *drone.Netrc
	Repo   *drone.Repo
	Stage  *drone.Stage
	System *drone.System
}

type execCommand struct {
	*Flags

	Source     string
	Include    []string
	Exclude    []string
	Privileged []string
	Networks   []string
	Volumes    map[string]string
	Environ    map[string]string
	Labels     map[string]string
	Secrets    map[string]string
	Resources  compiler.Resources
	Tmate      compiler.Tmate
	Clone      bool
	Config     string
	Pretty     bool
	Procs      int64
	Debug      bool
	Trace      bool
	Dump       bool
	PublicKey  string
	PrivateKey string
}

func mapOldToExecCommand(input *cli.Context) *execCommand {
	pipelineFile := input.Args().First()
	if pipelineFile == "" {
		pipelineFile = ".drone.yml"
	}
	return &execCommand{
		Flags: &Flags{
			Build: &drone.Build{
				Event:  input.String("event"),
				Ref:    input.String("ref"),
				Deploy: input.String("deploy-to"),
				Target: input.String("branch"),
			},
			Repo: &drone.Repo{
				Trusted: input.Bool("trusted"),
				Timeout: int64(input.Duration("timeout").Seconds()),
				Branch:  input.String("branch"),
				Name:    input.String("name"),
			},
			Stage: &drone.Stage{
				Name: input.String("pipeline"),
			},
			Netrc: &drone.Netrc{
				Machine:  input.String("netrc-machine"),
				Login:    input.String("netrc-username"),
				Password: input.String("netrc-password"),
			},
			System: &drone.System{
				Host: input.String("instance"),
			},
		},
		Source:     pipelineFile,
		Include:    input.StringSlice("include"),
		Exclude:    input.StringSlice("exclude"),
		Clone:      input.Bool("clone"),
		Networks:   input.StringSlice("network"),
		Environ:    readParams(input.String("env-file")),
		Volumes:    withVolumeSlice(input.StringSlice("volume")),
		Secrets:    readParams(input.String("secrets")),
		Config:     input.String("registry"),
		Privileged: input.StringSlice("privileged"),
		Pretty:     input.BoolT("pretty"),
	}
}

// WithVolumeSlice is a transform function that adds a set of global volumes to the container that are defined in --volume=host:container format.
func withVolumeSlice(volumes []string) (to map[string]string) {
	to = map[string]string{}
	for _, s := range volumes {
		parts := strings.Split(s, ":")
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		val := parts[1]
		to[key] = val
	}
	return to
}

// helper function reads secrets from a key-value file.
func readParams(path string) map[string]string {
	data, _ := godotenv.Read(path)
	return data
}
