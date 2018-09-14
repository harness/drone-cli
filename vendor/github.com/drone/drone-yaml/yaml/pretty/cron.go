package pretty

import "github.com/drone/drone-yaml/yaml"

// helper function pretty prints the cron resource.
func printCron(w writer, v *yaml.Cron) {
	w.WriteString("---")
	w.WriteTagValue("kind", v.Kind)
	w.WriteTagValue("name", v.Name)
	printSpec(w, v)
	w.WriteByte('\n')
	w.WriteByte('\n')
}

// helper function pretty prints the spec block.
func printSpec(w writer, v *yaml.Cron) {
	w.WriteTag("spec")

	w.IndentIncrease()
	w.WriteTagValue("schedule", v.Spec.Schedule)
	w.WriteTagValue("branch", v.Spec.Branch)
	if hasDeployment(v) {
		printDeploy(w, v)
	}
	w.IndentDecrease()
}

// helper function pretty prints the deploy block.
func printDeploy(w writer, v *yaml.Cron) {
	w.WriteTag("deployment")
	w.IndentIncrease()
	w.WriteTagValue("target", v.Spec.Deploy.Target)
	w.IndentDecrease()
}

// helper function returns true if the deployment
// object is empty.
func hasDeployment(v *yaml.Cron) bool {
	return v.Spec.Deploy.Target != ""
}
