package builder

import "io"

// ResultWriter represents the result from a build request.
type ResultWriter struct {
	writer   io.Writer
	exitCode int
}

// Write writes the build stdout and stderr to the result.
func (r *ResultWriter) Write(p []byte) (n int, err error) {
	return r.writer.Write(p)
}

// WriteExitCode writes the build exit status to the result.
func (r *ResultWriter) WriteExitCode(code int) {
	r.exitCode = code
}

// ExitCode returns the build exit status.
func (r *ResultWriter) ExitCode() int {
	return r.exitCode
}
