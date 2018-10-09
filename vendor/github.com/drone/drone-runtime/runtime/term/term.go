package term

import (
	"fmt"
	"io"
	"sync"

	"github.com/drone/drone-runtime/runtime"
)

const (
	linePlain  = "[%s:%d] %s"
	linePretty = "\033[%s[%s:%d]\033[0m %s"
)

// available terminal colors
var colors = []string{
	"32m", // green
	"33m", // yellow
	"34m", // blue
	"35m", // magenta
	"36m", // cyan
}

// WriteLineFunc defines a function responsible for writing indidual
// lines of log output.
type WriteLineFunc func(*runtime.State, *runtime.Line) error

// WriteLine writes log lines to io.Writer w in plain text format.
func WriteLine(w io.Writer) WriteLineFunc {
	return func(state *runtime.State, line *runtime.Line) error {
		fmt.Fprintf(w, linePlain, state.Step.Metadata.Name, line.Number, line.Message)
		return nil
	}
}

// WriteLinePretty writes pretty-printed log lines to io.Writer w.
func WriteLinePretty(w io.Writer) WriteLineFunc {
	var (
		mutex sync.Mutex
		steps = map[string]string{}
	)

	return func(state *runtime.State, line *runtime.Line) error {
		mutex.Lock()
		color, ok := steps[state.Step.Metadata.Name]
		mutex.Unlock()

		if !ok {
			color = colors[len(steps)%len(colors)]
			mutex.Lock()
			steps[state.Step.Metadata.Name] = color
			mutex.Unlock()
		}

		fmt.Fprintf(w, linePretty, color, state.Step.Metadata.Name, line.Number, line.Message)
		return nil
	}
}
