package runner

type Runner interface {
	Register(*Container) error
	Run(ResultWriter) error
	Cancel() error
}
