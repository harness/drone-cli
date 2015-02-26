package runner

import (
	"encoding/json"
	"io"

	"github.com/drone/drone-cli/config"
	"github.com/samalba/dockerclient"
)

type Response struct {
	Writer   io.Writer
	ExitCode int
}

func (r *Response) Write(p []byte) (n int, err error) {
	return r.Writer.Write(p)
}

func (r *Response) WriteExitCode(code int) {
	r.ExitCode = code
}

type Request struct {
	Clone  *Clone              `json:"clone"`
	Config *config.Config      `json:"config"`
	Client dockerclient.Client `json:"-"`
}

func (r *Request) Encode() string {
	out, _ := json.Marshal(r)
	return string(out)
}

type Clone struct {
	Origin  string   `json:"origin"`
	Remote  string   `json:"remote"`
	Branch  string   `json:"branch"`
	Sha     string   `json:"sha"`
	Ref     string   `json:"ref"`
	Dir     string   `json:"dir"`
	Netrc   *Netrc   `json:"netrc"`
	Keypair *Keypair `json:"keypair"`
}

type Netrc struct {
	Machine  string `json:"machine"`
	Login    string `json:"login"`
	Password string `json:"user"`
}

type Keypair struct {
	Public  string `json:"public"`
	Private string `json:"private"`
}
