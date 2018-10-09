package compiler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml/yaml"
)

func setupScriptWin(spec *engine.Spec, dst *engine.Step, src *yaml.Container) {
	var buf bytes.Buffer
	for _, command := range src.Commands {
		escaped := fmt.Sprintf("%q", command)
		escaped = strings.Replace(escaped, "$", `\$`, -1)
		buf.WriteString(fmt.Sprintf(
			traceScriptWin,
			escaped,
			command,
		))
	}
	script := fmt.Sprintf(
		buildScriptWin,
		buf.String(),
	)
	dst.Docker.Command = []string{"powershell", "-noprofile", "-noninteractive", "-command"}
	dst.Docker.Args = []string{"[System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String($Env:CI_SCRIPT)) | iex"}
	dst.Envs["CI_SCRIPT"] = base64.StdEncoding.EncodeToString([]byte(script))
	dst.Envs["HOME"] = "c:\\root"
	dst.Envs["SHELL"] = "powershell.exe"
}

// buildScriptWin is a helper script this is added to the build
// to prepare the environment and execute the build commands.
const buildScriptWin = `
$ErrorActionPreference = 'Stop';
&cmd /c "mkdir c:\root";
if ($Env:CI_NETRC_MACHINE) {
$netrc=[string]::Format("{0}\_netrc",$Env:HOME);
"machine $Env:CI_NETRC_MACHINE" >> $netrc;
"login $Env:CI_NETRC_USERNAME" >> $netrc;
"password $Env:CI_NETRC_PASSWORD" >> $netrc;
};
[Environment]::SetEnvironmentVariable("CI_NETRC_PASSWORD",$null);
[Environment]::SetEnvironmentVariable("CI_SCRIPT",$null);
[Environment]::SetEnvironmentVariable("DRONE_NETRC_USERNAME",$null);
[Environment]::SetEnvironmentVariable("DRONE_NETRC_PASSWORD",$null);
%s
`

// traceScriptWin is a helper script that is added to
// the build script to trace a command.
const traceScriptWin = `
Write-Output ('+ %s');
& %s; if ($LASTEXITCODE -ne 0) {exit $LASTEXITCODE}
`
