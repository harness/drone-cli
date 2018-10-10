package runtime

// Hook provides a set of hooks to run at various stages of
// runtime execution.
type Hook struct {
	// Before is called before all all steps are executed.
	Before func(*State) error

	// BeforeEach is called before each step is executed.
	BeforeEach func(*State) error

	// After is called after all steps are executed.
	After func(*State) error

	// AfterEach is called after each step is executed.
	AfterEach func(*State) error

	// GotLine is called when a line is logged.
	GotLine func(*State, *Line) error

	// GotLogs is called when the logs are completed.
	GotLogs func(*State, []*Line) error
}
