package builder

// Builder executes a set of build step handlers.
type Builder struct {
	handlers []Handler
}

// Handle registers the build step handler.
func (b *Builder) Handle(h Handler) {
	b.handlers = append(b.handlers, h)
}

// Build runs all build step handlers.
func (b *Builder) Build(rw *ResultWriter) (err error) {
	for _, h := range b.handlers {
		err = h.Build(rw)
		if err != nil || rw.exitCode != 0 {
			break
		}
	}
	return nil
}

// Cancel cancels all running build steps.
func (b *Builder) Cancel() {
	for _, h := range b.handlers {
		h.Cancel() // TODO use channel to signal cancel
	}
}
