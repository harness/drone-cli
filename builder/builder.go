package builder

type Builder struct {
	handlers []Handler
}

func (b *Builder) Register(handler Handler) {
	b.handlers = append(b.handlers, handler)
}

func (b *Builder) Build(rw ResponseWriter) error {
	for _, h := range b.handlers {
		err := h.Run(rw)
		if err != nil {
			return err
		}
		if rw.ExitCode() != 0 {
			break
		}
	}
	return nil
}

func (b *Builder) Cancel() {
	for _, h := range b.handlers {
		h.Cancel()
	}
}
