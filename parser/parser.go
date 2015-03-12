package parser

import (
	"io/ioutil"

	"github.com/drone/drone-cli/common"
	"gopkg.in/yaml.v2"
)

// Parse parses a build matrix and returns
// a list of build configurations for each axis.
func Parse(raw string) ([]*common.Config, error) {
	axis, err := ParseMatrix(raw)
	if err != nil {
		return nil, err
	}
	confs := []*common.Config{}

	// when no matrix values exist we should return
	// a single config value with an empty label.
	if len(axis) == 0 {
		conf, err := parse(raw)
		if err != nil {
			return nil, err
		}
		confs = append(confs, conf)
	}

	for _, ax := range axis {
		// inject the matrix values into the raw script
		injected := Inject(raw, ax)
		conf, err := parse(injected)
		if err != nil {
			return nil, err
		}
		conf.Axis = ax
		confs = append(confs, conf)
	}
	return confs, nil
}

// ParseFile parses a build matrix from a file and returns
// a list of build configurations for each axis.
func ParseFile(name string) ([]*common.Config, error) {
	raw, _ := ioutil.ReadFile(name)
	return Parse(string(raw))
}

// helper funtion to parse a yaml configuration file.
func parse(raw string) (*common.Config, error) {
	cfg := common.Config{}
	err := yaml.Unmarshal([]byte(raw), &cfg)
	return &cfg, err
}
