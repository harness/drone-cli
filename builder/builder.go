package builder

type Builder struct {
	handlers []Handler
}

// Handle adds a build step handler to be processed
// when running the build.
func (b *Builder) Handle(h Handler) {
	b.handlers = append(b.handlers, h)
}

// Build runs all build step handlers.
func (b *Builder) Build(r *Result) (err error) {
	for _, h := range b.handlers {
		err = h.Build(r)
		if err != nil || r.exitCode != 0 {
			break
		}
	}
	return nil
}

// Cancel cancels any running build processes and
// removes and build containers.
func (b *Builder) Cancel() {
	for _, h := range b.handlers {
		h.Cancel() // TODO use channel to signal cancel
	}
}

// type Builder struct {
// 	containers []*Container
// 	client     dockerclient.Client
// }
//
// // New creates a new builder
// func New(client dockerclient.Client) *Builder {
// 	return &Builder{client: client}
// }
//
// func (r *Builder) Add(c *Container) {
// 	r.containers = append(r.containers, c)
// }
//
// func (r *Builder) Cancel() {
// 	for _, c := range r.containers {
// 		// TODO cancel using environment
// 		if c != nil {
// 			// TODO remove
// 		}
// 	}
// }
//
// func (b *Builder) Build(res *Result) error {
// 	for _, c := range b.containers {
// 		err := b.run(c, res)
// 		if err != nil {
// 			return err
// 		}
// 		if res.ExitCode() != 0 {
// 			return nil
// 		}
// 	}
// 	return nil
// }
//
// // helper function to run a single build container.
// func (b *Builder) run(c *Container, res *Result) error {
// 	// create container
// 	//    err: failed creating
// 	// start continer
// 	//    err: failed starting container
// 	if c.Detached == false {
// 		return nil
// 	}
// 	// get log reader
// 	//    err: failed getting logs
//
// 	// copy logs to writer
// 	//    err: failed copying logs
//
// 	// write response
// 	return nil
// }

/*
func NewBuilder(build *Build) *Builder {
	b := New(build.Client)
	for _, step := range build.Config.Compose {
		b.Add(fromCompose(build, &step))
	}
	b.Add(fromSetup(build, &build.Config.Build))
	b.Add(fromPlugin(build, &build.Config.Clone))
	b.Add(fromBuild(build, &build.Config.Build))

	return b
}

func NewDeployer(build *Build) *Builder {
	b := New(build.Client)
	for _, step := range build.Config.Publish {
		b.Add(fromPlugin(build, &step))
	}
	for _, step := range build.Config.Deploy {
		b.Add(fromPlugin(build, &step))
	}
	return b
}

func NewNotifier(build *Build) *Builder {
	b := New(build.Client)
	for _, step := range build.Config.Notify {
		b.Add(fromPlugin(build, &step))
	}
	return b
}
*/
