package runner

/*
import (
	"io"

	"github.com/samalba/dockerclient"
)

type Runner struct {
	client     dockerclient.Client
	containers []*Container
}

func (r *Runner) Run(req *Request, resp ResponseWriter) error {

	return nil
}

Setup()

Clone()
Build()
Deploy()
Notify()

Teardown()

func (r *Runner) Logs(w io.Writer) {
	for _, c := range r.containers {
		if c.Detached {
			continue
		}
		r, err := c.Logs()
		if err != nil {
			continue
		}
		io.Copy(w, r)
	}
}

func (r *Runner) Kill() {
	for _, c := range r.containers {
		c.Stop()
		c.Kill()
		c.Remove()
	}
}
*/
