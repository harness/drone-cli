package docker

import (
	"archive/tar"
	"bytes"
	"time"

	"github.com/drone/drone-runtime/engine"
)

// helper function creates a tarfile that can be uploaded
// into the Docker container.
func createTarfile(file *engine.File, mount *engine.FileMount) []byte {
	w := new(bytes.Buffer)
	t := tar.NewWriter(w)
	h := &tar.Header{
		Name:    mount.Path,
		Mode:    mount.Mode,
		Size:    int64(len(file.Data)),
		ModTime: time.Now(),
	}
	t.WriteHeader(h)
	t.Write(file.Data)
	return w.Bytes()
}
