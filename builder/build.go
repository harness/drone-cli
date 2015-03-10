package builder

import (
	"encoding/json"

	"github.com/drone/drone-cli/common"
	"github.com/samalba/dockerclient"
)

// Build represents a build request.
type Build struct {
	Repo   *common.Repo
	Commit *common.Repo
	Config *common.Config
	Clone  *common.Clone

	Client dockerclient.Client
}

// BuildPayload represents the payload of a plugin
// that is serialized and sent to the plugin in JSON
// format via stdin or arg[1].
type BuildPayload struct {
	Repo   *common.Repo  `json:"repo"`
	Commit *common.Repo  `json:"commit"`
	Clone  *common.Clone `json:"clone"`

	Config map[string]interface{} `json:"vargs"`
}

// Encode encodes the payload in JSON format.
func (b *BuildPayload) Encode() string {
	out, _ := json.Marshal(b)
	return string(out)
}
