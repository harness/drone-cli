package pretty

import (
	"github.com/drone/drone-yaml/yaml"
)

// helper function to pretty print the pipeline resource.
func printPipeline(w writer, v *yaml.Pipeline) {
	w.WriteString("---")
	w.WriteTagValue("kind", v.Kind)
	w.WriteTagValue("type", v.Type)
	w.WriteTagValue("name", v.Name)
	w.WriteByte('\n')

	if !isPlatformEmpty(v.Platform) {
		printPlatform(w, v.Platform)
	} else {
		printPlatformDefault(w)
	}
	if !isCloneEmpty(v.Clone) {
		printClone(w, v.Clone)
	}
	if !isWorkspaceEmpty(v.Workspace) {
		printWorkspace(w, v.Workspace)
	}

	if len(v.Steps) > 0 {
		w.WriteTag("steps")
		for _, step := range v.Steps {
			seq := new(indexWriter)
			seq.writer = w
			seq.IndentIncrease()
			printContainer(seq, step)
			seq.IndentDecrease()
		}
	}

	if len(v.Services) > 0 {
		w.WriteTag("services")
		for _, step := range v.Services {
			seq := new(indexWriter)
			seq.writer = w
			seq.IndentIncrease()
			printContainer(seq, step)
			seq.IndentDecrease()
		}
	}

	if len(v.Volumes) != 0 {
		printVolumes(w, v.Volumes)
		w.WriteByte('\n')
	}

	if !isConditionsEmpty(v.Trigger) {
		printConditions(w, "trigger", v.Trigger)
		w.WriteByte('\n')
	}

	if len(v.DependsOn) > 0 {
		printDependsOn(w, v.DependsOn)
		w.WriteByte('\n')
	}

	w.WriteByte('\n')
}

// helper function pretty prints the clone block.
func printClone(w writer, v yaml.Clone) {
	w.WriteTag("clone")
	w.IndentIncrease()
	w.WriteTagValue("depth", v.Depth)
	w.WriteTagValue("disable", v.Disable)
	w.WriteByte('\n')
	w.IndentDecrease()
}

// helper function pretty prints the conditions mapping.
func printConditions(w writer, name string, v yaml.Conditions) {
	w.WriteTag(name)
	w.IndentIncrease()
	if !isConditionEmpty(v.Branch) {
		printCondition(w, "branch", v.Branch)
	}
	if !isConditionEmpty(v.Event) {
		printCondition(w, "event", v.Event)
	}
	if !isConditionEmpty(v.Instance) {
		printCondition(w, "instance", v.Instance)
	}
	if !isConditionEmpty(v.Paths) {
		printCondition(w, "paths", v.Paths)
	}
	if !isConditionEmpty(v.Ref) {
		printCondition(w, "ref", v.Ref)
	}
	if !isConditionEmpty(v.Repo) {
		printCondition(w, "repo", v.Repo)
	}
	if !isConditionEmpty(v.Status) {
		printCondition(w, "status", v.Status)
	}
	if !isConditionEmpty(v.Target) {
		printCondition(w, "target", v.Target)
	}
	w.IndentDecrease()
}

// helper function pretty prints a condition mapping.
func printCondition(w writer, k string, v yaml.Condition) {
	w.WriteTag(k)
	if len(v.Include) != 0 && len(v.Exclude) == 0 {
		w.WriteByte('\n')
		w.Indent()
		writeValue(w, v.Include)
	}
	if len(v.Include) != 0 && len(v.Exclude) != 0 {
		w.IndentIncrease()
		w.WriteTagValue("include", v.Include)
		w.IndentDecrease()
	}
	if len(v.Exclude) != 0 {
		w.IndentIncrease()
		w.WriteTagValue("exclude", v.Include)
		w.IndentDecrease()
	}
}

// helper function pretty prints the target platform.
func printPlatform(w writer, v yaml.Platform) {
	w.WriteTag("platform")
	w.IndentIncrease()
	w.WriteTagValue("os", v.OS)
	w.WriteTagValue("arch", v.Arch)
	w.WriteTagValue("variant", v.Variant)
	w.WriteTagValue("version", v.Version)
	w.WriteByte('\n')
	w.IndentDecrease()
}

// helper function prints default platform values.
// Including target platform is considered a best-practive.
func printPlatformDefault(w writer) {
	w.WriteTag("platform")
	w.IndentIncrease()
	w.WriteTagValue("os", "linux")
	w.WriteTagValue("arch", "amd64")
	w.WriteByte('\n')
	w.IndentDecrease()
}

// helper function pretty prints the volume sequence.
func printVolumes(w writer, v []*yaml.Volume) {
	w.WriteTag("volumes")
	for _, v := range v {
		s := new(indexWriter)
		s.writer = w
		s.IndentIncrease()

		s.WriteTagValue("name", v.Name)
		if v := v.EmptyDir; v != nil {
			s.WriteTag("temp")
			s.IndentIncrease()
			s.WriteTagValue("medium", v.Medium)
			s.WriteTagValue("size_limit", v.SizeLimit)
			s.IndentDecrease()
		}

		if v := v.HostPath; v != nil {
			s.WriteTag("host")
			s.IndentIncrease()
			s.WriteTagValue("path", v.Path)
			s.IndentDecrease()
		}

		s.IndentDecrease()
	}
}

// helper function pretty prints the workspace block.
func printWorkspace(w writer, v yaml.Workspace) {
	w.WriteTag("workspace")
	w.IndentIncrease()
	w.WriteTagValue("base", v.Base)
	w.WriteTagValue("path", v.Path)
	w.WriteByte('\n')
	w.IndentDecrease()
}

// helper function returns true if the workspace
// object is empty.
func isWorkspaceEmpty(v yaml.Workspace) bool {
	return v.Path == "" && v.Base == ""
}

// helper function returns true if the platform
// object is empty.
func isPlatformEmpty(v yaml.Platform) bool {
	return v.OS == "" &&
		v.Arch == "" &&
		v.Variant == "" &&
		v.Version == ""
}

// helper function returns true if the clone
// object is empty.
func isCloneEmpty(v yaml.Clone) bool {
	return v.Depth == 0 &&
		v.Disable == false
}

// helper function returns true if the conditions
// object is empty.
func isConditionsEmpty(v yaml.Conditions) bool {
	return isConditionEmpty(v.Branch) &&
		isConditionEmpty(v.Event) &&
		isConditionEmpty(v.Instance) &&
		isConditionEmpty(v.Paths) &&
		isConditionEmpty(v.Ref) &&
		isConditionEmpty(v.Repo) &&
		isConditionEmpty(v.Status) &&
		isConditionEmpty(v.Target)
}

// helper function returns true if the condition
// object is empty.
func isConditionEmpty(v yaml.Condition) bool {
	return len(v.Exclude) == 0 && len(v.Include) == 0
}
