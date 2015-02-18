package config

import (
	"fmt"

	"gopkg.in/yaml.v1"
)

const (
	LimitAxis  = 10
	LimitPerms = 25
)

func Parse(raw string) (*Config, error) {
	config := Config{}
	err := yaml.Unmarshal([]byte(raw), &config)
	return &config, err
}

func ParseMatrix(raw string) ([]*Config, error) {
	var matrix []*Config

	config, err := Parse(raw)
	if err != nil {
		return matrix, err
	}

	// if not a matrix build return an array
	// with just the single axis.
	if len(config.Matrix) == 0 {
		matrix = append(matrix, config)
		return matrix, nil
	}

	// calculate number of permutations and
	// extract the list of keys.
	var perm int
	var keys []string
	for k, v := range config.Matrix {
		perm *= len(v)
		if perm == 0 {
			perm = len(v)
		}
		keys = append(keys, k)
	}

	// for each axis calculate the values the uniqe
	// set of values that should be used.
	for p := 0; p < perm; p++ {
		axis := map[string]string{}
		decr := perm
		for i, key := range keys {
			vals := config.Matrix[key]
			decr = decr / len(vals)
			item := p / decr % len(vals)
			axis[key] = vals[item]

			// enforce a maximum number of axis
			// in the build matrix.
			if i > LimitAxis {
				break
			}
		}

		config, err = Parse(Inject(raw, axis))
		if err != nil {
			return nil, err
		}
		matrix = append(matrix, config)

		// each axis value should also be added
		// as an environment variable
		for key, val := range axis {
			env := fmt.Sprintf("%s=%s", key, val)
			config.Env = append(config.Env, env)
		}

		// enforce a maximum number of permutations
		// in the build matrix.
		if p > LimitPerms {
			break
		}
	}

	return matrix, nil
}
