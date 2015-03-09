package docker

type Container struct {
	ID          string
	Image       string
	Pull        bool
	Detach      bool
	Privileged  bool
	Entrypoint  []string
	Cmd         []string
	Env         []string
	Volumes     []string
	WorkingDir  string
	NetworkMode string
}

type State struct {
	Running  bool
	Pid      int
	ExitCode int
	Started  int64
	Finished int64
}
