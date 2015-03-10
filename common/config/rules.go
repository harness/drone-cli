package config

import (
	"strings"

	"github.com/drone/drone-cli/common"
)

// Rule applies a check or transformation rule to
// the build configuration.
type Rule func(*common.Config)

// addSetup is a transformer that adds a default
// setup step if none exists.
func addSetup(c *common.Config) {
	c.Setup = &common.Step{}
	c.Setup.Image = "plugins/drone-setup"
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
	c.Deploy = map[string]*common.Step{}
}

// rmDeploy is a transformer that removes all
// publish steps.
func rmDeploy(c *common.Config) {
	c.Publish = map[string]*common.Step{}
}

// rmPrivileged is a transformer that ensures every
// step is executed in non-privileged mode.
func rmPrivileged(c *common.Config) {
	c.Clone.Privileged = false
	c.Build.Privileged = false
	for _, step := range c.Publish {
		step.Privileged = false
	}
	for _, step := range c.Deploy {
		step.Privileged = false
	}
	for _, step := range c.Notify {
		step.Privileged = false
	}
}

// rmVolumes is a transformer that ensures every
// step is executed without volumes.
func rmVolumes(c *common.Config) {
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
}

// rmNetwork is a transformer that ensures every
// step is executed with default bridge networking.
func rmNetwork(c *common.Config) {
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
