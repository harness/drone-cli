{
  new(
    image,
    commands=[],
    detach=false,
    environment=[],
    group='',
    secrets=[],
    when=null,
    pull=true,
  ):: {
    image: image,
    [if commands != [] then 'commands']: commands,
    [if detach == true then 'detach']: true,
    [if environment != [] then 'environment']: environment,
    [if group != '' then 'group']: group,
    [if secrets != [] then 'secrets']: secrets,
    [if when != null then 'when']: when,
    [if pull == true then 'pull']: pull,

    with(params):: self + params,
  },
}