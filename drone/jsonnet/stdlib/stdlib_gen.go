package stdlib

import jsonnet "github.com/google/go-jsonnet"

var files = map[string]jsonnet.Contents{
	"drone.libsonnet":       jsonnet.MakeContents("{\n  v1:: import 'v1/v1.libsonnet',\n}"),
	"v1/metadata.libsonnet": jsonnet.MakeContents("{\n  new(\n    name\n  ):: {\n    name: name,\n  },\n}"),
	"v1/platform.libsonnet": jsonnet.MakeContents("{\n  new(\n    os='linux',\n    arch='amd64',\n    variant=null,\n    kernel=null\n  ):: {\n    os: os,\n    arch: arch,\n    [if kernel != null then 'kernel']: kernel,\n    [if variant != null then 'variant']: variant,\n  },\n}"),
	"v1/secrets.libsonnet":  jsonnet.MakeContents("{\n  new(\n    name,\n    external=null,\n    secretbox=null\n  ):: {\n    [name]: {\n      [if external != null then 'external']: external,\n      [if secretbox != null then 'secretbox']: secretbox,\n    },\n  },\n}"),
	"v1/step.libsonnet":     jsonnet.MakeContents("{\n  new(\n    image,\n    commands=[],\n    detach=false,\n    environment=[],\n    group='',\n    secrets=[],\n    when=null,\n    pull=true,\n  ):: {\n    image: image,\n    [if commands != [] then 'commands']: commands,\n    [if detach == true then 'detach']: true,\n    [if environment != [] then 'environment']: environment,\n    [if group != '' then 'group']: group,\n    [if secrets != [] then 'secrets']: secrets,\n    [if when != null then 'when']: when,\n    [if pull == true then 'pull']: pull,\n\n    with(params):: self + params,\n  },\n}"),
	"v1/v1.libsonnet":       jsonnet.MakeContents("{\n  metadata:: import 'metadata.libsonnet',\n  platform:: import 'platform.libsonnet',\n  secrets:: import 'secrets.libsonnet',\n  step:: import 'step.libsonnet',\n}\n"),
}
