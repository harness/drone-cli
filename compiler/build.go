package compiler

import (
	"io"
	"time"

	"github.com/drone/drone-cli/common"
	"github.com/samalba/dockerclient"
)

// B is a type passed to build nodes. B implements an io.Writer
// and will accumulate build output during execution.
type B struct {
	Repo   *common.Repo
	Commit *common.Commit
	Config *common.Config
	Clone  *common.Clone

	client dockerclient.Client

	writer io.Writer

	exitCode int

	start    time.Time // Time build started
	duration time.Duration
	timerOn  bool

	containers []string
}

// NewB returns a new Build context.
func NewB(client dockerclient.Client, w io.Writer) *B {
	return &B{
		client: client,
		writer: w,
	}
}

// Run creates and runs a Docker container.
func (b *B) Run(conf *dockerclient.ContainerConfig) (string, error) {
	name, err := b.client.CreateContainer(conf, "")
	if err != nil {
		// on error try to pull the Docker image.
		// note that this may not be the cause of
		// the error, but we'll try just in case.
		b.client.PullImage(conf.Image, nil)

		// then try to re-create
		name, err = b.client.CreateContainer(conf, "")
		if err != nil {
			return name, err
		}
	}
	b.containers = append(b.containers, name)
	err = b.client.StartContainer(name, &conf.HostConfig)
	if err != nil {
		return name, err
	}

	return name, nil
}

// Inspect inspects the running Docker container and returns
// the contianer runtime information and state.
func (b *B) Inspect(name string) (*dockerclient.ContainerInfo, error) {
	return b.client.InspectContainer(name)
}

// Remove stops and removes the named Docker container.
func (b *B) Remove(name string) {
	b.client.StopContainer(name, 5)
	b.client.KillContainer(name, "SIGKILL")
	b.client.RemoveContainer(name, true, true)
}

// RemoveAll stops and removes all Docker containers that were
// created and started during the build process.
func (b *B) RemoveAll() {
	for _, name := range b.containers {
		b.Remove(name)
	}
}

// Logs returns an io.ReadCloser for reading the build stream of
// the named Docker container.
func (b *B) Logs(name string) (io.ReadCloser, error) {
	opts := dockerclient.LogOptions{
		Follow:     true,
		Stderr:     true,
		Stdout:     true,
		Timestamps: true,
	}
	return b.client.ContainerLogs(name, &opts)
}

// StartTimer starts timing a build. This function is called automatically
// before a build starts, but it can also used to resume timing after
// a call to StopTimer.
func (b *B) StartTimer() {
	if !b.timerOn {
		b.start = time.Now()
		b.timerOn = true
	}
}

// StopTimer stops timing a build. This can be used to pause the timer
// while performing complex initialization that you don't want to measure.
func (b *B) StopTimer() {
	if b.timerOn {
		b.duration += time.Now().Sub(b.start)
		b.timerOn = false
	}
}

// Write writes the build stdout and stderr to the result.
func (b *B) Write(p []byte) (n int, err error) {
	return b.writer.Write(p)
}

// Exit writes the function as having failed but continues execution.
func (b *B) Exit(code int) {
	if code != 0 { // never override non-zero exit
		b.exitCode = code
	}
}

// ExitCode reports the build exit code. A non-zero value indicates
// the build exited with errors.
func (b *B) ExitCode() int {
	return b.exitCode
}
