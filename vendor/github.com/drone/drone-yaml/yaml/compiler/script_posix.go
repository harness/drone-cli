package compiler

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml/yaml"
	"github.com/drone/drone-yaml/yaml/compiler/internal/rand"
)

func setupScript(spec *engine.Spec, dst *engine.Step, src *yaml.Container) {
	var buf bytes.Buffer
	for _, command := range src.Commands {
		escaped := fmt.Sprintf("%q", command)
		escaped = strings.Replace(escaped, "$", `\$`, -1)
		buf.WriteString(fmt.Sprintf(
			traceScript,
			escaped,
			command,
		))
	}
	script := fmt.Sprintf(
		buildScript,
		buf.String(),
	)
	spec.Files = append(spec.Files, &engine.File{
		Metadata: engine.Metadata{
			UID:       rand.String(),
			Namespace: spec.Metadata.Namespace,
			Name:      src.Name,
		},
		Data: []byte(script),
	})
	dst.Files = append(dst.Files, &engine.FileMount{
		Name: src.Name,
		Path: "/usr/drone/bin/init",
		Mode: 0777,
	})

	dst.Docker.Command = []string{"/bin/sh"}
	dst.Docker.Args = []string{"/usr/drone/bin/init"}
}

// buildScript is a helper script this is added to the build
// to prepare the environment and execute the build commands.
const buildScript = `
if [ -n "$CI_NETRC_MACHINE" ]; then
cat <<EOF > $HOME/.netrc
machine $CI_NETRC_MACHINE
login $CI_NETRC_USERNAME
password $CI_NETRC_PASSWORD
EOF
chmod 0600 $HOME/.netrc
fi
unset CI_NETRC_USERNAME
unset CI_NETRC_PASSWORD
unset DRONE_NETRC_USERNAME
unset DRONE_NETRC_PASSWORD
set -e
%s
`

// traceScript is a helper script that is added to
// the build script to trace a command.
const traceScript = `
echo + %s
%s
`
