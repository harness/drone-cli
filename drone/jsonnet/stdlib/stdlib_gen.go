package stdlib

import jsonnet "github.com/google/go-jsonnet"

var files = map[string]jsonnet.Contents{
	"drone.libsonnet":    jsonnet.MakeContents("{\n  metadata:: import 'metadata.libsonnet',\n  platform:: import 'platform.libsonnet',\n  secrets:: import 'secrets.libsonnet',\n  step:: import 'step.libsonnet',\n}\n"),
	"metadata.libsonnet": jsonnet.MakeContents("{\n  new(name):: {\n    name: name,\n  },\n}\n"),
	"platform.libsonnet": jsonnet.MakeContents("{\n  new(\n    os='linux',\n    arch='amd64',\n    arm=null,\n    version=null\n  ):: {\n    os: os,\n    arch: arch,\n    [if arm != null then 'arm']: arm,\n    [if version != null then 'version']: version,\n  },\n}\n"),
	"secrets.libsonnet":  jsonnet.MakeContents("{\n  new(\n    name,\n    external=null,\n    secretbox=null\n  ):: {\n    [name]: {\n      [if external != null then 'external']: external,\n      [if secretbox != null then 'secretbox']: secretbox,\n    },\n  },\n}\n"),
	"step.libsonnet":     jsonnet.MakeContents("{\n  new(\n    name,\n    image,\n    commands=[],\n    detach=false,\n    environment=[],\n    group='',\n    secrets=[],\n    when=null,\n    pull=true,\n  ):: {\n    [name]: {\n      image: image,\n      [if commands != [] then 'commands']: commands,\n      [if detach == true then 'detach']: true,\n      [if environment != [] then 'environment']: environment,\n      [if group != '' then 'group']: group,\n      [if secrets != [] then 'secrets']: secrets,\n      [if when != null then 'when']: when,\n      [if pull == true then 'pull']: pull,\n    },\n  },\n\n  docker:: import './steps/docker.libsonnet',\n  githubRelease:: import './steps/githubRelease.libsonnet',\n  gpgsign:: import './steps/gpgsign.libsonnet',\n  manifest:: import './steps/manifest.libsonnet',\n  webhook:: import './steps/webhook.libsonnet',\n}\n"),
}
