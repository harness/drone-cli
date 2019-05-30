// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Non-Commercial License
// that can be found in the LICENSE file.

package yaml

type (
	// Constraints defines a set of runtime constraints.
	Constraints struct {
		Ref         Constraint
		Repo        Constraint
		Instance    Constraint
		Environment Constraint
		Event       Constraint
		Branch      Constraint
		Status      Constraint
		Matrix      map[string]string
	}

	// Constraint defines a runtime constraint.
	Constraint struct {
		Include []string
		Exclude []string
	}

	// ConstraintMap defines a runtime constraint map.
	ConstraintMap struct {
		Include map[string]string
		Exclude map[string]string
	}
)

// UnmarshalYAML unmarshals the constraint.
func (c *Constraint) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 = struct {
		Include StringSlice
		Exclude StringSlice
	}{}

	var out2 StringSlice

	unmarshal(&out1)
	unmarshal(&out2)

	c.Exclude = out1.Exclude
	c.Include = append(
		out1.Include,
		out2...,
	)
	return nil
}

// UnmarshalYAML unmarshals the constraint map.
func (c *ConstraintMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	out1 := struct {
		Include map[string]string
		Exclude map[string]string
	}{
		Include: map[string]string{},
		Exclude: map[string]string{},
	}

	out2 := map[string]string{}

	unmarshal(&out1)
	unmarshal(&out2)

	c.Include = out1.Include
	c.Exclude = out1.Exclude
	for k, v := range out2 {
		c.Include[k] = v
	}
	return nil
}
