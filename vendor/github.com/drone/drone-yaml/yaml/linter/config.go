// Copyright 2019 Drone IO, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package linter

import (
	"errors"

	"github.com/drone/drone-yaml/yaml"
)

// ErrDuplicatePipelineName is returned when two Pipeline
// resources have the same name.
var ErrDuplicatePipelineName = errors.New("linter: duplicate pipeline names")

// ErrMissingPipelineDependency is returned when a Pipeline
// defines dependencies that are invlid or unknown.
var ErrMissingPipelineDependency = errors.New("linter: invalid or unknown pipeline dependency")

// ErrCyclicalPipelineDependency is returned when a Pipeline
// defines a cyclical dependency, which would result in an
// infinite execution loop.
var ErrCyclicalPipelineDependency = errors.New("linter: cyclical pipeline dependency detected")

// ErrPipelineSelfDependency is returned when a Pipeline
// defines a dependency on itself.
var ErrPipelineSelfDependency = errors.New("linter: pipeline cannot have a dependency on itself")

// Manifest performs lint operations for a manifest.
func Manifest(manifest *yaml.Manifest, trusted bool) error {
	return checkPipelines(manifest, trusted)
}

func checkPipelines(manifest *yaml.Manifest, trusted bool) error {
	names := map[string]struct{}{}
	for _, resource := range manifest.Resources {
		switch v := resource.(type) {
		case *yaml.Pipeline:
			_, ok := names[v.Name]
			if ok {
				return ErrDuplicatePipelineName
			}
			names[v.Name] = struct{}{}
			err := checkPipelineDeps(v, names)
			if err != nil {
				return err
			}
			err = checkPipeline(v, trusted)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func checkPipelineDeps(pipeline *yaml.Pipeline, deps map[string]struct{}) error {
	for _, dep := range pipeline.DependsOn {
		_, ok := deps[dep]
		if !ok {
			return ErrMissingPipelineDependency
		}
		if pipeline.Name == dep {
			return ErrPipelineSelfDependency
		}
	}
	return nil
}
