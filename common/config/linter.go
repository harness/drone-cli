package config

import (
	"fmt"

	"github.com/drone/drone-cli/common"
)

// lintRule defines a function that runs lint
// checks against a Yaml Config file. If the rule
// fails it should return an error message.
type lintRule func(*common.Config) error

var lintRules = []lintRule{
	expectBuild, expectImage, expectCommand,
}

// Lint runs all lint rules against the Yaml Config.
func Lint(c *common.Config) error {
	for _, rule := range lintRules {
		err := rule(c)
		if err != nil {
			return err
		}
	}
	return nil
}

// lint rule that fails when no build is defined
func expectBuild(c *common.Config) error {
	if c.Build == nil {
		return fmt.Errorf("Yaml must define a build section")
	}
	return nil
}

// lint rule that fails when no build image is defined
func expectImage(c *common.Config) error {
	if len(c.Build.Image) == 0 {
		return fmt.Errorf("Yaml must define a build image")
	}
	return nil
}

// lint rule that fails when no build commands are defined
func expectCommand(c *common.Config) error {
	if c.Build.Config == nil || c.Build.Config["commands"] == nil {
		return fmt.Errorf("Yaml must define build commands")
	}
	return nil
}
