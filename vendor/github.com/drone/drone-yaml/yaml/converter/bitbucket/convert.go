// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Non-Commercial License
// that can be found in the LICENSE file.

package bitbucket

import (
	"bytes"
	"fmt"

	droneyaml "github.com/drone/drone-yaml/yaml"
	"github.com/drone/drone-yaml/yaml/pretty"

	"gopkg.in/yaml.v2"
)

// Convert converts the yaml configuration file from
// the legacy format to the 1.0+ format.
func Convert(b []byte, ref string) ([]byte, error) {
	config := new(Config)
	err := yaml.Unmarshal(b, config)
	if err != nil {
		return nil, err
	}

	// TODO (bradrydzewski) to correctly choose
	// the pipeline we need to pass the branch
	// and ref.
	stage := config.Pipeline(ref)

	pipeline := &droneyaml.Pipeline{}
	pipeline.Name = "default"
	pipeline.Kind = "pipeline"

	//
	// clone
	//

	pipeline.Clone.Depth = config.Clone.Depth

	//
	// steps
	//

	for i, from := range stage.Steps {
		to := toContainer(from)
		// defaults to the global image if the
		// step does not define an image.
		if to.Image == "" {
			to.Image = config.Image
		}
		if to.Name == "" {
			to.Name = fmt.Sprintf("step_%d", i)
		}
		pipeline.Steps = append(pipeline.Steps, to)
	}

	//
	// services
	//

	for name, from := range config.Definitions.Services {
		to := toContainer(from)
		to.Name = name
		pipeline.Services = append(pipeline.Services, to)
	}

	//
	// wrap the pipeline in the manifest
	//

	manifest := &droneyaml.Manifest{}
	manifest.Resources = append(manifest.Resources, pipeline)

	buf := new(bytes.Buffer)
	pretty.Print(buf, manifest)
	return buf.Bytes(), nil
}

func toContainer(from *Step) *droneyaml.Container {
	return &droneyaml.Container{
		Name:     from.Name,
		Image:    from.Image,
		Commands: from.Script,
	}
}
