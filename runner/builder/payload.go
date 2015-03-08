package builder

import (
	"encoding/json"

	"github.com/drone/drone-cli/common"
)

type Payload struct {
	Clone  *common.Clone  `json:"clone"`
	Commit *common.Commit `json:"commit"`
	Repo   *common.Repo   `json:"repo"`
	User   *common.User   `json:"user"`
	Vargs  interface{}    `json:"vargs"`
}

func NewPayload(r *Request, args interface{}) *Payload {
	return &Payload{
		Clone:  r.Clone,
		Commit: r.Commit,
		Repo:   r.Repo,
		User:   r.User,
		Vargs:  args,
	}
}

func (p *Payload) Encode() string {
	out, _ := json.Marshal(p)
	return string(out)
}
