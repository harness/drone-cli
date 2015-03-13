package compiler

import "github.com/drone/drone-cli/common"

// Builder represents a build execution tree.
type Builder struct {
	Build  Node
	Deploy Node
	Notify Node
}

// Run runs the build, deploy and notify nodes
// in the build tree.
func (b *Builder) Run() error {
	var err error
	err = b.RunBuild()
	if err != nil {
		return err
	}
	err = b.RunDeploy()
	if err != nil {
		return err
	}
	return b.RunNotify()
}

// RunBuild runs only the build node.
func (b *Builder) RunBuild() error {
	return nil
}

// RunDeploy runs only the deploy node.
func (b *Builder) RunDeploy() error {
	return nil
}

// RunNotify runs on the notify node.
func (b *Builder) RunNotify() error {
	return nil
}

// Load loads a build configuration file.
func Load(conf *common.Config) *Builder {
	var (
		builds  []Node
		deploys []Node
		notifys []Node
	)

	for _, step := range conf.Compose {
		builds = append(builds, &batchNode{step}) // compose
	}
	builds = append(builds, &batchNode{}) // setup
	builds = append(builds, &batchNode{}) // clone
	builds = append(builds, &batchNode{}) // build

	for _, step := range conf.Publish {
		deploys = append(deploys, &batchNode{step}) // publish
	}
	for _, step := range conf.Deploy {
		deploys = append(deploys, &batchNode{step}) // deploy
	}
	for _, step := range conf.Notify {
		notifys = append(notifys, &batchNode{step}) // notify
	}
	return &Builder{
		serialNode(builds),
		serialNode(deploys),
		serialNode(notifys),
	}
}
