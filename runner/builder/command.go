package builder

type Command struct {
	ID          string
	Image       string
	Pull        bool
	Privileged  bool
	Entrypoint  []string
	Cmd         []string
	Env         []string
	Volumes     []string
	NetworkMode string

	Defer  bool
	Detach bool
}

// run
// run background
// run post
// run regardless (ie notify)
