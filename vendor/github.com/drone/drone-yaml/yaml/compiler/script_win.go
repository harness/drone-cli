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
	dst.Envs["SHELL"] = "powershell.exe"
}

// buildScriptWin is a helper script this is added to the build
// to prepare the environment and execute the build commands.
const buildScriptWin = `
if ($Env:CI_NETRC_MACHINE) {
@"
machine $Env:CI_NETRC_MACHINE
login $Env:CI_NETRC_USERNAME
password $Env:CI_NETRC_PASSWORD
"@ > (Join-Path $Env:USERPROFILE '_netrc');
}
[Environment]::SetEnvironmentVariable("CI_NETRC_USERNAME", $null);
[Environment]::SetEnvironmentVariable("CI_NETRC_PASSWORD", $null);
[Environment]::SetEnvironmentVariable("DRONE_NETRC_USERNAME", $null);
[Environment]::SetEnvironmentVariable("DRONE_NETRC_PASSWORD", $null);
[Environment]::SetEnvironmentVariable("CI_SCRIPT", $null);
$ErrorActionPreference = 'Stop';
%s
`

// traceScriptWin is a helper script that is added to
// the build script to trace a command.
const traceScriptWin = `
Write-Output ('+ %s');
& %s; if ($LASTEXITCODE -ne 0) {exit $LASTEXITCODE}
`
