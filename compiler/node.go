package compiler

import (
	"sync"

	"github.com/drone/drone-cli/common"
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
	// host *dockerclient.HostConfig
	// conf *dockerclient.ContainerConfig
}

func (n *batchNode) Run(s *State, rw *ResultWriter) error {
	// host := toHostConfig(n.step)
	// conf := toContainerConfig(n.step)
	// if step.Config != nil {
	// 	conf.Cmd = toCommand(s, n.step)
	// }
	// name, err := s.Run(conf, host)
	// if err != nil {
	// 	return nil
	// }
	//
	// rc, err := s.Logs(name)
	// if err != nil {
	// 	return err
	// }
	// io.Copy(rw, rc)
	//
	// info, err := c.Info(name)
	// if err != nil {
	// 	return err
	// }
	// rw.WriteExitCode(info.State.ExitCode)
	switch {
	case n.step.Condition == nil:
	case n.step.Condition.MatchBranch(s.Commit.Branch) == false:
		return nil
	case n.step.Condition.MatchOwner(s.Repo.Owner) == false:
		return nil
	}

	return nil
}

// serviceNode runs a container, blocking, writes output, uses config section
type serviceNode struct {
	step *common.Step
}

func (n *serviceNode) Run(s *State, rw *ResultWriter) error {
	// host := toHostConfig(n.step)
	// conf := toContainerConfig(n.step)

	// _, err := s.Run(conf, host)
	// return err
	return nil
}
