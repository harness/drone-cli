package engine

import (
	"bytes"
	"encoding/json"
)

// PullPolicy defines the container image pull policy.
type PullPolicy int

// PullPolicy enumeration.
const (
	PullDefault PullPolicy = iota
	PullAlways
	PullIfNotExists
	PullNever
)

func (p PullPolicy) String() string {
	return pullPolicyID[p]
}

var pullPolicyID = map[PullPolicy]string{
	PullDefault:     "default",
	PullAlways:      "always",
	PullIfNotExists: "if-not-exists",
	PullNever:       "never",
}

var pullPolicyName = map[string]PullPolicy{
	"":              PullDefault,
	"default":       PullDefault,
	"always":        PullAlways,
	"if-not-exists": PullIfNotExists,
	"never":         PullNever,
}

// MarshalJSON marshals the string representation of the
// pull type to JSON.
func (p *PullPolicy) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(pullPolicyID[*p])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals the json representation of the
// pull type from a string value.
func (p *PullPolicy) UnmarshalJSON(b []byte) error {
	// unmarshal as string
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	// lookup value
	*p = pullPolicyName[s]
	return nil
}

// RunPolicy defines the policy for starting containers
// based on the point-in-time pass or fail state of
// the pipeline.
type RunPolicy int

// RunPolicy enumeration.
const (
	RunOnSuccess RunPolicy = iota
	RunOnFailure
	RunAlways
	RunNever
)

func (r RunPolicy) String() string {
	return runPolicyID[r]
}

var runPolicyID = map[RunPolicy]string{
	RunOnSuccess: "on-success",
	RunOnFailure: "on-failure",
	RunAlways:    "always",
	RunNever:     "never",
}

var runPolicyName = map[string]RunPolicy{
	"":           RunOnSuccess,
	"on-success": RunOnSuccess,
	"on-failure": RunOnFailure,
	"always":     RunAlways,
	"never":      RunNever,
}

// MarshalJSON marshals the string representation of the
// run type to JSON.
func (r *RunPolicy) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(runPolicyID[*r])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals the json representation of the
// run type from a string value.
func (r *RunPolicy) UnmarshalJSON(b []byte) error {
	// unmarshal as string
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	// lookup value
	*r = runPolicyName[s]
	return nil
}
