package builder

import "io"

// Result represents the result from a build request.
type Result struct {
	writer   io.Writer
	exitCode int
}

// Write writes the build stdout and stderr to the result.
func (r *Result) Write(p []byte) (n int, err error) {
	return r.writer.Write(p)
}

// WriteExitCode writes the build exit status to the result.
func (r *Result) WriteExitCode(code int) {
	r.exitCode = code
}

// ExitCode returns the build exit status.
func (r *Result) ExitCode() int {
	return r.exitCode
}
