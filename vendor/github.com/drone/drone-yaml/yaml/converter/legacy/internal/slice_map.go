// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Non-Commercial License
// that can be found in the LICENSE file.

package yaml

import "strings"

// SliceMap represents a slice or map of key pairs.
type SliceMap struct {
	Map map[string]string
}

// UnmarshalYAML implements custom Yaml unmarshaling.
func (s *SliceMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	s.Map = map[string]string{}
	err := unmarshal(&s.Map)
	if err == nil {
		return nil
	}

	var slice []string
	err = unmarshal(&slice)
	if err != nil {
		return err
	}
	for _, v := range slice {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			val := parts[1]
			s.Map[key] = val
		}
	}
	return nil
}
