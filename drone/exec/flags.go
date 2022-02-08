// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package exec

import (
	"github.com/drone-runners/drone-runner-docker/engine/compiler"
	"github.com/drone/drone-go/drone"
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

func mapOldToExecCommand(input *cli.Context) (returnVal *execCommand) {
	returnVal = &execCommand{
		Flags: &Flags{
			Build: &drone.Build{
				Event:  input.String("event"),
				Ref:    input.String("ref"),
				Deploy: input.String("deploy-to"),
			},
			Repo: &drone.Repo{
				Trusted: input.Bool("trusted"),
				Timeout: int64(input.Int("timeout")),
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
		Source:     input.Args().First(),
		Include:    input.StringSlice("include"),
		Exclude:    input.StringSlice("exclude"),
		Clone:      input.Bool("clone"),
		Networks:   input.StringSlice("network"),
		Privileged: input.StringSlice("privileged"),
	}

	return returnVal
}
