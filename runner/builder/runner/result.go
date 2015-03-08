package runner

import "io"

// A ResultWriter interface is used by the builder to
// construct a build repsonse.
type ResultWriter interface {
	// Write writes the build stdout and stderr to the response.
	Write([]byte) (int, error)

	// WriteExitCode writes the build exit status to the response.
	// Explicit calls to WriteExitCode should be used to signal
	// build failures or errors.
	WriteExitCode(int)

	// ExitCode returns the build exit status. If the exit status
	// has not been set this will return 0.
	ExitCode() int
}

// Result represents the response resulting from
// build execution.
type Result struct {
	writer   io.Writer
	exitCode int
}

func (r *Result) Write(p []byte) (n int, err error) {
	// TODO(brydzews) this is the perfect spot to parse
	// the Docker log format and convert to a standard
	// plain text stream.
	return r.writer.Write(p)
}

func (r *Result) WriteExitCode(code int) {
	r.exitCode = code
}

func (r *Result) ExitCode() int {
	return r.exitCode
}
