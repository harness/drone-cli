package runner

import "io"

// A ResponseWriter interface is used by the builder to
// construct a build repsonse.
type ResponseWriter interface {
	// Write writes the build stdout and stderr to the response.
	Write([]byte) (int, error)

	// WriteExitCode writes the build exit status to the response.
	// Explicit calls to WriteExitCode should be used to signal
	// build failures or errors.
	WriteExitCode(int)
}

// Response represents the response resulting from
// build execution.
type Response struct {
	Writer   io.Writer
	ExitCode int
}

func (r *Response) Write(p []byte) (n int, err error) {
	// TODO(brydzews) this is the perfect spot to convert
	// parse the Docker log format and convert to a standard
	// plain text stream.
	return r.Writer.Write(p)
}

func (r *Response) WriteExitCode(code int) {
	r.ExitCode = code
}
