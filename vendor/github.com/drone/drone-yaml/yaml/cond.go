package yaml

import filepath "github.com/bmatcuk/doublestar"

// Conditions defines a group of conditions.
type Conditions struct {
	Ref      Condition
	Repo     Condition
	Instance Condition
	Target   Condition
	Event    Condition
	Branch   Condition
	Status   Condition
	Paths    Condition
}

// Condition defines a runtime condition.
type Condition struct {
	Include []string `yaml:"include,omitempty" json:"include,omitempty"`
	Exclude []string `yaml:"exclude,omitempty" json:"exclude,omitempty"`
}

// Match returns true if the string matches the include
// patterns and does not match any of the exclude patterns.
func (c *Condition) Match(v string) bool {
	if c.Excludes(v) {
		return false
	}
	if c.Includes(v) {
		return true
	}
	if len(c.Include) == 0 {
		return true
	}
	return false
}

// Includes returns true if the string matches the include
// patterns.
func (c *Condition) Includes(v string) bool {
	for _, pattern := range c.Include {
		if ok, _ := filepath.Match(pattern, v); ok {
			return true
		}
	}
	return false
}

// Excludes returns true if the string matches the exclude
// patterns.
func (c *Condition) Excludes(v string) bool {
	for _, pattern := range c.Exclude {
		if ok, _ := filepath.Match(pattern, v); ok {
			return true
		}
	}
	return false
}

// UnmarshalYAML implements yml unmarhsaling.
func (c *Condition) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 []string
	var out2 = struct {
		Include []string
		Exclude []string
	}{}

	unmarshal(&out1)
	unmarshal(&out2)

	c.Exclude = out2.Exclude
	c.Include = append(
		out2.Include,
		out1...,
	)
	return nil
}
