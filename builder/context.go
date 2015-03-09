package builder

import (
	"github.com/drone/drone-cli/engine"
)

type Context struct {
	Engine   engine.Engine
	Request  *Request
	Response *Response
}
