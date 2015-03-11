package config

import (
	"strings"

	"github.com/drone/drone-cli/common"
)

// transformRule applies a check or transformation rule
// to the build configuration.
type transformRule func(*common.Config)

// Transform executes the default transformers that
// ensure the minimal Yaml configuration is in place
// and correctly configured.
func Transform(c *common.Config) {
	addSetup(c)
	addClone(c)
	normalizeBuild(c)
	normalizeImages(c)
	normalizeDockerPlugin(c)
}

// TransformSafe executes all transformers that remove
// privileged options from the Yaml.
func TransformSafe(c *common.Config) {
	rmPrivileged(c)
	rmVolumes(c)
	rmNetwork(c)
}

// TransformBuild executes all transformers that remove
// non-build and non-notfiy steps from the Yaml.
func TransformBuild(c *common.Config) {
	rmPublish(c)
	rmDeploy(c)
}

// addSetup is a transformer that adds a default
// setup step if none exists.
func addSetup(c *common.Config) {
	c.Setup = &common.Step{}
	c.Setup.Image = "plugins/drone-build"
	c.Setup.Config = c.Build.Config
}

// addClone is a transformer that adds a default
// clone step if none exists.
func addClone(c *common.Config) {
	if c.Clone == nil {
		c.Clone = &common.Step{}
	}
	if len(c.Clone.Image) == 0 {
		c.Clone.Image = "plugins/drone-git"
	}
}

// normalizeBuild is a transformer that removes the
// build configuration vargs. They should have
// already been transferred to the Setup step.
func normalizeBuild(c *common.Config) {
	c.Build.Config = nil
}

// normalizeImages is a transformer that ensures every
// step has an image and uses a fully-qualified
// image name.
func normalizeImages(c *common.Config) {
	c.Setup.Image = imageName(c.Setup.Image)
	c.Clone.Image = imageName(c.Clone.Image)
	c.Build.Image = imageName(c.Build.Image)
	for name, step := range c.Publish {
		step.Image = imageNameDefault(step.Image, name)
	}
	for name, step := range c.Deploy {
		step.Image = imageNameDefault(step.Image, name)
	}
	for name, step := range c.Notify {
		step.Image = imageNameDefault(step.Image, name)
	}
}

// normalizeDockerPlugin is a transformer that ensures the
// official Docker plugin can runs in privileged mode. It
// will disable volumes and network mode for added protection.
func normalizeDockerPlugin(c *common.Config) {
	for _, step := range c.Publish {
		if step.Image == "plugins/drone-docker" {
			step.Privileged = true
			step.Volumes = []string{}
			step.NetworkMode = ""
			break
		}
	}
}

// rmPublish is a transformer that removes all
// publish steps.
func rmPublish(c *common.Config) {
	c.Publish = map[string]*common.Step{}
}

// rmDeploy is a transformer that removes all
// publish steps.
func rmDeploy(c *common.Config) {
	c.Deploy = map[string]*common.Step{}
}

// rmNotify is a transformer that removes all
// notify steps.
func rmNotify(c *common.Config) {
	c.Notify = map[string]*common.Step{}
}

// rmPrivileged is a transformer that ensures every
// step is executed in non-privileged mode.
func rmPrivileged(c *common.Config) {
	c.Setup.Privileged = false
	c.Clone.Privileged = false
	c.Build.Privileged = false
	for _, step := range c.Publish {
		if step.Image == "plugins/drone-docker" {
			continue // the official docker plugin is the only exception here
		}
		step.Privileged = false
	}
	for _, step := range c.Deploy {
		step.Privileged = false
	}
	for _, step := range c.Notify {
		step.Privileged = false
	}
	for _, step := range c.Compose {
		step.Privileged = false
	}
}

// rmVolumes is a transformer that ensures every
// step is executed without volumes.
func rmVolumes(c *common.Config) {
	c.Setup.Volumes = []string{}
	c.Clone.Volumes = []string{}
	c.Build.Volumes = []string{}
	for _, step := range c.Publish {
		step.Volumes = []string{}
	}
	for _, step := range c.Deploy {
		step.Volumes = []string{}
	}
	for _, step := range c.Notify {
		step.Volumes = []string{}
	}
	for _, step := range c.Compose {
		step.Volumes = []string{}
	}
}

// rmNetwork is a transformer that ensures every
// step is executed with default bridge networking.
func rmNetwork(c *common.Config) {
	c.Setup.NetworkMode = ""
	c.Clone.NetworkMode = ""
	c.Build.NetworkMode = ""
	for _, step := range c.Publish {
		step.NetworkMode = ""
	}
	for _, step := range c.Deploy {
		step.NetworkMode = ""
	}
	for _, step := range c.Notify {
		step.NetworkMode = ""
	}
	for _, step := range c.Compose {
		step.NetworkMode = ""
	}
}

// imageName is a helper function that resolves the
// image name. When using official drone plugins it
// is possible to use an alias name. This converts to
// the fully qualified name.
func imageName(name string) string {
	if strings.Contains(name, "/") {
		return name
	}
	name = strings.Replace(name, "_", "-", -1)
	name = "plugins/drone-" + name
	return name
}

// imageNameDefault is a helper function that resolves
// the image name. If the image name is blank the
// default name is used instead.
func imageNameDefault(name, defaultName string) string {
	if len(name) == 0 {
		name = defaultName
	}
	return imageName(name)
}
