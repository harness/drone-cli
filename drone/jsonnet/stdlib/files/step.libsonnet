{
  new(
    name,
    image,
    commands=[],
    detach=false,
    environment=[],
    group='',
    secrets=[],
    when=null,
    pull=true,
  ):: {
    [name]: {
      image: image,
      [if commands != [] then 'commands']: commands,
      [if detach == true then 'detach']: true,
      [if environment != [] then 'environment']: environment,
      [if group != '' then 'group']: group,
      [if secrets != [] then 'secrets']: secrets,
      [if when != null then 'when']: when,
      [if pull == true then 'pull']: pull,
    },
  },

  docker:: import './steps/docker.libsonnet',
  githubRelease:: import './steps/githubRelease.libsonnet',
  gpgsign:: import './steps/gpgsign.libsonnet',
  manifest:: import './steps/manifest.libsonnet',
  webhook:: import './steps/webhook.libsonnet',
}
