// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Non-Commercial License
// that can be found in the LICENSE file.

package yaml

import (
	"bytes"

	"github.com/vinzenz/yaml"
)

type temporary struct {
	Attributes map[string]interface{} `yaml:",inline"`
	Pipeline   yaml.MapSlice          `yaml:"pipeline"`
}

// this is a helper function that expands merge keys
func expandMergeKeys(b []byte) ([]byte, error) {
	v := new(temporary)
	if err := yaml.Unmarshal(b, v); err != nil {
		return b, err
	}
	o, err := yaml.Marshal(v)
	if err != nil {
		return b, err
	}
	return o, nil
}

func hasMergeKeys(b []byte) bool {
	return bytes.Contains(b, []byte("<<:"))
}
