// Copyright 2019 Drone IO, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package runtime

import (
	"strings"
	"time"

	"github.com/drone/drone-runtime/engine"
)

// Line represents a line in the container logs.
type Line struct {
	Number    int    `json:"pos,omitempty"`
	Message   string `json:"out,omitempty"`
	Timestamp int64  `json:"time,omitempty"`
}

type lineWriter struct {
	num   int
	now   time.Time
	rep   *strings.Replacer
	state *State
	lines []*Line
	size  int
	limit int
}

func newWriter(state *State) *lineWriter {
	w := &lineWriter{}
	w.num = 0
	w.now = time.Now().UTC()
	w.state = state
	w.rep = newReplacer(state.config.Secrets)
	w.limit = 5242880 // 5MB max log size
	return w
}

func (w *lineWriter) Write(p []byte) (n int, err error) {
	// if the maximum log size has been exceeded, the
	// log entry is silently ignored.
	if w.size >= w.limit {
		return len(p), nil
	}

	out := string(p)
	if w.rep != nil {
		out = w.rep.Replace(out)
	}

	parts := []string{out}

	// kubernetes buffers the output and may combine
	// multiple lines into a single block of output.
	// Split into multiple lines.
	//
	// note that docker output always inclines a line
	// feed marker. This needs to be accounted for when
	// splitting the output into multiple lines.
	if strings.Contains(strings.TrimSuffix(out, "\n"), "\n") {
		parts = strings.SplitAfter(out, "\n")
	}

	for _, part := range parts {
		line := &Line{
			Number:    w.num,
			Message:   part,
			Timestamp: int64(time.Since(w.now).Seconds()),
		}

		if w.state.hook.GotLine != nil {
			w.state.hook.GotLine(w.state, line)
		}
		w.size = w.size + len(part)
		w.num++

		w.lines = append(w.lines, line)
	}

	// if the write exceeds the maximum output we should
	// write a single line to the end of the logs that
	// indicates the output is being truncated.
	if w.size >= w.limit {
		w.lines = append(w.lines, &Line{
			Number:    w.num,
			Message:   "warning: maximum output exceeded",
			Timestamp: int64(time.Since(w.now).Seconds()),
		})
	}

	return len(p), nil
}

func newReplacer(secrets []*engine.Secret) *strings.Replacer {
	var oldnew []string
	for _, secret := range secrets {
		oldnew = append(oldnew, secret.Data)
		oldnew = append(oldnew, "********")
	}
	if len(oldnew) == 0 {
		return nil
	}
	return strings.NewReplacer(oldnew...)
}
