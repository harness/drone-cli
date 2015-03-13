package compiler

import (
	"io"

	"github.com/drone/drone-cli/common"
	"github.com/samalba/dockerclient"
)

// State represents the execution state of the
// running build.
type State struct {
	Repo   *common.Repo
	Commit *common.Repo
	Config *common.Config
	Clone  *common.Clone

	Client dockerclient.Client
}

func (*State) Run() (string, error) {
	return "", nil
}

func (*State) Remove(string) error {
	return nil
}

func (*State) RemoveAll() error {
	return nil
}

func (*State) Logs(string) (io.ReadCloser, error) {
	return nil, nil
}
