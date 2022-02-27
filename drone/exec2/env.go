package exec2

import (
	"os"
	"strings"

	"github.com/urfave/cli"
)

func getEnv(c *cli.Context) map[string]string {
	env := prefixedEnviron(
		os.Environ(),
	)
	if c.IsSet("branch") {
		v := c.String("branch")
		env["DRONE_BRANCH"] = v
		env["DRONE_COMMIT_BRANCH"] = v
		env["DRONE_TARGET_BRANCH"] = v
	}
	if c.IsSet("event") {
		v := c.String("event")
		env["DRONE_EVENT"] = v
	}
	if c.IsSet("instance") {
		v := c.String("instance")
		env["DRONE_SYSTEM_HOST"] = v
		env["DRONE_SYSTEM_HOSTNAME"] = v
	}
	if c.IsSet("ref") {
		v := c.String("ref")
		env["DRONE_COMMIT_REF"] = v
	}
	if c.IsSet("sha") {
		v := c.String("sha")
		env["DRONE_COMMIT_SHA"] = v
	}
	if c.IsSet("repo") {
		v := c.String("repo")
		env["DRONE_REPO"] = v
	}
	if c.IsSet("deploy-to") {
		v := c.String("deploy-to")
		env["DRONE_DEPLOY_TO"] = v
	}
	return env
}

// helper function returns all environment variables
// prefixed with DRONE_.
func prefixedEnviron(environ []string) map[string]string {
	envs := map[string]string{}
	for _, env := range environ {
		if !strings.HasPrefix(env, "DRONE_") {
			continue
		}
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		val := parts[1]
		envs[key] = val
	}
	return envs
}
