package compiler

import (
	"sync"

	"github.com/drone/drone-cli/common"
	"github.com/samalba/dockerclient"
)

// Node is an element in the build execution tree.
type Node interface {
	Run(*State, *ResultWriter) error
}

// parallelNode runs a set of build nodes in parallel.
type parallelNode []Node

func (n parallelNode) Run(s *State, rw *ResultWriter) error {
	var wg sync.WaitGroup
	for _, node := range n {
		wg.Add(1)

		go func(node Node) {
			defer wg.Done()
			node.Run(s, rw)
		}(node)
	}
	wg.Wait()
	return nil
}

// serialNode runs a set of build nodes in sequential order.
type serialNode []Node

func (n serialNode) Run(s *State, rw *ResultWriter) error {
	for _, node := range n {
		err := node.Run(s, rw)
		if err != nil {
			return err
		}
		if rw.ExitCode() != 0 {
			return nil
		}
	}
	return nil
}

// batchNode runs a container and blocks until complete.
type batchNode struct {
	step *common.Step
}

func (n *batchNode) Run(s *State, rw *ResultWriter) error {
	host := toHostConfig(n.step)
	conf := toContainerConfig(n.step)
	if n.step.Config != nil {
		conf.Cmd = toCommand(s, n.step)
		conf.Entrypoint = []string{}
	}

	return nil
}

// serviceNode runs a container, blocking, writes output, uses config section
type serviceNode struct {
	step *common.Step
}

func (n *serviceNode) Run(s *State, rw *ResultWriter) error {
	host := toHostConfig(n.step)
	conf := toContainerConfig(n.step)

	return nil
}

// cloneNode runs a clone container, blocking, uses build section
type cloneNode struct {
	step *common.Step
}

func (n *cloneNode) Run(s *State, rw *ResultWriter) error {
	return nil
}

// buildNode runs a build container, discards the step.config section
type buildNode struct {
	host *dockerclient.HostConfig
	conf *dockerclient.ContainerConfig
}

func newBuildNode(step *common.Step) *buildNode {
	host := toHostConfig(step)
	conf := toContainerConfig(step)
	conf.Entrypoint = []string{"/bin/bash"}
	conf.Cmd = []string{"/drone/bin/build.sh"}
	return &buildNode{host, conf}
}

func (n *buildNode) Run(s *State, rw *ResultWriter) error {
	return nil
}

// setupNode container, discards the step.config section
type setupNode struct {
	step *common.Step
}

func (n *setupNode) Run(s *State, rw *ResultWriter) error {
	return nil
}
