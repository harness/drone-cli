package builder

import (
	"io"
)

type ResponseWriter interface {
	// Write writes the build stdout and stderr to the response.
	Write([]byte) (int, error)

	// WriteExitCode writes the build exit status to the response.
	WriteExitCode(int)
}

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
