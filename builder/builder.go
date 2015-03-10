package builder

import "sync"

// Builder executes a set of build step handlers.
type Builder struct {
	sync.Mutex
	handlers []Handler
}

// Handle registers the build step handler.
func (b *Builder) Handle(h Handler) {
	b.Lock()
	defer b.Unlock()
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
	b.Lock()
	defer b.Unlock()

	for _, h := range b.handlers {
		h.Cancel() // TODO use channel to signal cancel
	}
}
