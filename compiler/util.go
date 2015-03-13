package compiler

import (
	"encoding/json"

	"github.com/drone/drone-cli/common"
	"github.com/samalba/dockerclient"
)

// helper function to encode the build step to
// a json string. Primarily used for plugins, which
// expect a json encoded string in stdin or arg[1].
func toCommand(state *State, step *common.Step) []string {
	p := payload{
		state.Repo,
		state.Commit,
		state.Clone,
		step.Config,
	}
	return []string{p.Encode()}
}

// helper function that converts a build step to
// a hostConfig for use with the dockerclient
func toHostConfig(step *common.Step) *dockerclient.HostConfig {
	return &dockerclient.HostConfig{
		Privileged:  step.Privileged,
		NetworkMode: step.NetworkMode,
	}
}

// helper function that converts the build step to
// a containerConfig for use with the dockerclient
func toContainerConfig(step *common.Step) *dockerclient.ContainerConfig {
	config := &dockerclient.ContainerConfig{
		Image:      step.Image,
		Env:        step.Environment,
		Cmd:        step.Command,
		Entrypoint: step.Entrypoint,
		WorkingDir: step.WorkingDir,
	}

	if len(step.Volumes) != 0 {
		config.Volumes = map[string]struct{}{}
		for _, path := range step.Volumes {
			config.Volumes[path] = struct{}{}
		}
	}

	return config
}

// payload represents the payload of a plugin
// that is serialized and sent to the plugin in JSON
// format via stdin or arg[1].
type payload struct {
	Repo   *common.Repo  `json:"repo"`
	Commit *common.Repo  `json:"commit"`
	Clone  *common.Clone `json:"clone"`

	Config map[string]interface{} `json:"vargs"`
}

// Encode encodes the payload in JSON format.
func (p *payload) Encode() string {
	out, _ := json.Marshal(p)
	return string(out)
}
