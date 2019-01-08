package yaml

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type (
	// Containers represents an ordered list of containers.
	Containers struct {
		Containers []*Container
	}

	// Container represents a Docker container.
	Container struct {
		Command     StringSlice            `yaml:"command,omitempty"`
		Commands    StringSlice            `yaml:"commands,omitempty"`
		Detached    bool                   `yaml:"detach,omitempty"`
		Devices     []string               `yaml:"devices,omitempty"`
		ErrIgnore   bool                   `yaml:"allow_failure,omitempty"`
		Tmpfs       []string               `yaml:"tmpfs,omitempty"`
		DNS         StringSlice            `yaml:"dns,omitempty"`
		DNSSearch   StringSlice            `yaml:"dns_search,omitempty"`
		Entrypoint  StringSlice            `yaml:"entrypoint,omitempty"`
		Environment SliceMap               `yaml:"environment,omitempty"`
		ExtraHosts  []string               `yaml:"extra_hosts,omitempty"`
		Image       string                 `yaml:"image,omitempty"`
		Name        string                 `yaml:"name,omitempty"`
		Privileged  bool                   `yaml:"privileged,omitempty"`
		Pull        bool                   `yaml:"pull,omitempty"`
		Shell       string                 `yaml:"shell,omitempty"`
		Volumes     []*Volume              `yaml:"volumes,omitempty"`
		Secrets     Secrets                `yaml:"secrets,omitempty"`
		Constraints Constraints            `yaml:"when,omitempty"`
		Vargs       map[string]interface{} `yaml:",inline"`
	}
)

// UnmarshalYAML implements the Unmarshaller interface.
func (c *Containers) UnmarshalYAML(unmarshal func(interface{}) error) error {
	slice := yaml.MapSlice{}
	if err := unmarshal(&slice); err != nil {
		return err
	}

	for _, s := range slice {
		container := Container{}
		out, _ := yaml.Marshal(s.Value)

		if err := yaml.Unmarshal(out, &container); err != nil {
			return err
		}
		container.Name = fmt.Sprintf("%v", s.Key)
		c.Containers = append(c.Containers, &container)
	}
	return nil
}
