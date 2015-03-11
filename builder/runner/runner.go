package runner

import (
	"github.com/drone/drone-cli/builder"
)

// Builder is a convenience function that creates a build
// runner for build steps.
func Builder(build *builder.Build) *builder.Builder {
	b := builder.Builder{}
	for _, step := range build.Config.Compose {
		b.Handle(builder.Service(build, step))
	}
	b.Handle(builder.Batch(build, build.Config.Setup))
	b.Handle(builder.Batch(build, build.Config.Clone))
	b.Handle(builder.Batch(build, build.Config.Build))
	return &b
}

// Deployer is a convenience function that creates a build
// runner for publish and deploy steps.
func Deployer(build *builder.Build) *builder.Builder {
	b := builder.Builder{}
	for _, step := range build.Config.Publish {
		b.Handle(builder.Batch(build, step))
	}
	for _, step := range build.Config.Deploy {
		b.Handle(builder.Batch(build, step))
	}
	return &b
}

// Notifier is a convenience function that creates a build runner
// for notification steps.
func Notifier(build *builder.Build) *builder.Builder {
	b := builder.Builder{}
	for _, step := range build.Config.Notify {
		b.Handle(builder.Batch(build, step))
	}
	return &b
}
