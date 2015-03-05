package runner

import (
	"github.com/drone/drone-cli/common"
	"github.com/samalba/dockerclient"
)

// A Request represents a build request received by
// a build runner.
type Request struct {
	Clone  *common.Clone  `json:"clone"`
	Commit *common.Commit `json:"commit"`
	Repo   *common.Repo   `json:"repo"`
	User   *common.User   `json:"user"`

	// Config specifies the build configuration and execution
	// instructions to use when exeucting a build.
	Config *common.Config `json:"-"`

	// Client specifies the Docker client to use when executing
	// a build.
	Client dockerclient.Client `json:"-"`
}
