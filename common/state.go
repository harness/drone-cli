package common

type State struct {
	Running  bool
	Pid      int
	ExitCode int
	Started  int64
	Finished int64
}
