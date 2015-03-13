package parser

// ResultWriter represents the result from a build request.
type ResultWriter interface {

	// Write writes the build stdout and stderr to the result.
	Write([]byte) (int, error)

	// WriteExitCode writes the build exit status to the result.
	WriteExitCode(int)

	// ExitCode returns the build exit status.
	ExitCode() int
}
