package config

import "github.com/drone/drone-cli/common"

// Filter is a transform function that removes certain
// build steps that don't meet the configuration condition
// criteria (ie doesn't match branch).
func Filter(c *common.Config, repo *common.Repo, commit *common.Commit) {
	filterMap(c.Publish, repo, commit)
	filterMap(c.Deploy, repo, commit)
	filterMap(c.Notify, repo, commit)
}

// helper function that removes all steps from the given map
// that don't match the build conditions.
func filterMap(steps map[string]*common.Step, r *common.Repo, c *common.Commit) {
	for k, step := range steps {
		if isMatch(step, r, c) == false {
			delete(steps, k)
		}
	}
}

// helper function that returns false if the step does not
// match the build conditions.
func isMatch(step *common.Step, r *common.Repo, c *common.Commit) bool {
	switch {
	case step.Condition == nil:
		return true
	case step.Condition.MatchBranch(c.Branch):
		return false
	case step.Condition.MatchOwner(r.Owner):
		return false
	default:
		return true
	}
}
