// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Non-Commercial License
// that can be found in the LICENSE file.

package gitlab

import (
	"github.com/drone/drone-yaml/yaml/converter/internal"
)

type (
	// Config defines the pipeline configuration.
	Config struct {
		// Image specifies the Docker image with
		// which we run your builds.
		Image Image

		// Stages is used to group steps into stages,
		// where each stage is executed sequentially.
		Stages []string

		// Services is used to define a set of services
		// that should be started and linked to each
		// step in the pipeline.
		Services []*Image

		// Variables is used to customize execution,
		// such as the clone strategy.
		Variables map[string]string

		// Before contains the list of bash commands
		// that are executed in sequence before the
		// first job.
		Before internal.StringSlice `yaml:"before_script"`

		// After contains the list of bash commands
		// that are executed in sequence after the
		// last job.
		After internal.StringSlice `yaml:"after_script"`

		// Jobs is used to define individual units
		// of execution that make up a stage.
		Jobs map[string]*Job `yaml:",inline"`
	}

	// Job defines a build execution unit.
	Job struct {
		// Name of the pipeline step.
		Name string

		// Stage is the name of the stage.
		Stage string

		// Image specifies the Docker image with
		// which we run your builds.
		Image Image

		// Script contains the list of bash commands
		// that are executed in sequence.
		Script internal.StringSlice

		// Before contains the list of bash commands
		// that are executed in sequence before the
		// primary script.
		Before internal.StringSlice `yaml:"before_script"`

		// After contains the list of bash commands
		// that are executed in sequence after the
		// primary script.
		After internal.StringSlice `yaml:"after_script"`

		// Services defines a set of services linked
		// to the job.
		Services []*Image

		// Only defines the names of branches and tags
		// for which the job will run.
		Only internal.StringSlice

		// Except defines the names of branches and tags
		// for which the job will not run.
		Except internal.StringSlice

		// Variables is used to customize execution,
		// such as the clone strategy.
		Variables map[string]string

		// Allow job to fail. Failed job doesn’t contribute
		// to commit status
		AllowFailure bool

		// Define when to run job. Can be on_success, on_failure,
		// always or manual
		When internal.StringSlice
	}

	// Image defines a Docker image.
	Image struct {
		Name       string
		Entrypoint []string
		Command    []string
		Alias      string
	}
)

// UnmarshalYAML implements custom parsing for an Image.
func (i *Image) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var name string
	err := unmarshal(&name)
	if err == nil {
		i.Name = name
		return nil
	}
	data := struct {
		Name       string
		Entrypoint internal.StringSlice
		Command    internal.StringSlice
		Alias      string
	}{}
	err = unmarshal(&data)
	if err != nil {
		return err
	}
	i.Name = data.Name
	i.Entrypoint = data.Entrypoint
	i.Command = data.Command
	i.Alias = data.Alias
	return nil
}
