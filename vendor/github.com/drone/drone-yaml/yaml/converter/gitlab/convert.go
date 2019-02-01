package gitlab

import (
	"bytes"
	"strings"

	droneyaml "github.com/drone/drone-yaml/yaml"
	"github.com/drone/drone-yaml/yaml/compiler/image"
	"github.com/drone/drone-yaml/yaml/pretty"

	"gopkg.in/yaml.v2"
)

// Convert converts the yaml configuration file from
// the legacy format to the 1.0+ format.
func Convert(b []byte) ([]byte, error) {
	config := new(Config)
	err := yaml.Unmarshal(b, config)
	if err != nil {
		return nil, err
	}
	manifest := &droneyaml.Manifest{}

	// if no stages are defined, we create a single,
	// default stage that will be used for all jobs.
	if len(config.Stages) == 0 {
		for name, job := range config.Jobs {
			config.Stages = append(config.Stages, name)
			job.Stage = name
		}
	}

	// create a new pipeline for each stage.
	var prevstage string
	for _, stage := range config.Stages {
		pipeline := &droneyaml.Pipeline{}
		pipeline.Name = stage
		pipeline.Kind = droneyaml.KindPipeline
		manifest.Resources = append(manifest.Resources, pipeline)
		for name, job := range config.Jobs {
			if job.Stage != stage {
				continue
			}
			cmds := []string(config.Before)
			cmds = append(cmds, []string(job.Before)...)
			cmds = append(cmds, []string(job.Script)...)
			cmds = append(cmds, []string(job.After)...)
			cmds = append(cmds, []string(config.After)...)

			step := &droneyaml.Container{
				Name:       name,
				Image:      job.Image.Name,
				Command:    job.Image.Command,
				Entrypoint: job.Image.Entrypoint,
				Commands:   cmds,
			}

			if job.AllowFailure {
				step.Failure = "ignore"
			}

			if step.Image == "" {
				step.Image = config.Image.Name
			}
			// TODO: handle Services
			// TODO: handle Only
			// TODO: handle Except
			// TODO: handle Variables
			// TODO: handle When

			pipeline.Steps = append(pipeline.Steps, step)
		}

		for _, step := range config.Services {
			step := &droneyaml.Container{
				Name:       step.Alias,
				Image:      step.Name,
				Command:    step.Command,
				Entrypoint: step.Entrypoint,
			}
			if step.Name == "" {
				step.Name = serviceSlug(step.Image)
			}
			pipeline.Services = append(pipeline.Services, step)
		}

		if prevstage != "" {
			pipeline.DependsOn = []string{prevstage}
		}
		prevstage = stage
	}

	buf := new(bytes.Buffer)
	pretty.Print(buf, manifest)
	return buf.Bytes(), nil
}

func serviceSlug(s string) string {
	s = image.Trim(s)
	s = strings.Replace(s, "/", "__", -1)
	return s
}
