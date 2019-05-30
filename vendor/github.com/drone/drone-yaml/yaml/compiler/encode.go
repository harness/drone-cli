// Copyright 2019 Drone IO, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package compiler

import (
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/buildkite/yaml"
	json "github.com/ghodss/yaml"
)

// helper funciton encodes an interface value as a string.
// this function assumes all types were unmarshaled by the
// yaml.v2 library. The yaml.v2 package only supports a
// subset of primative types.
func encode(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case []byte:
		return base64.StdEncoding.EncodeToString(v)
	case []interface{}:
		return encodeSlice(v)
	default:
		return encodeMap(v)
	}
}

// helper function encodes a parameter in map format.
func encodeMap(v interface{}) string {
	yml, _ := yaml.Marshal(v)
	out, _ := json.YAMLToJSON(yml)
	return string(out)
}

// helper function encodes a parameter in slice format.
func encodeSlice(v interface{}) string {
	out, _ := yaml.Marshal(v)

	in := []string{}
	err := yaml.Unmarshal(out, &in)
	if err == nil {
		return strings.Join(in, ",")
	}
	out, _ = json.YAMLToJSON(out)
	return string(out)
}
