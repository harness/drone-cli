package runner

import "github.com/drone/drone-cli/common"

// A Request represents a build request received by
// a build runner.
type Request struct {
	Clone  *common.Clone
	Commit *common.Commit
	Repo   *common.Repo
	User   *common.User

	// Config specifies the build configuration and execution
	// instructions to use when exeucting a build.
	Config *common.Config
}
