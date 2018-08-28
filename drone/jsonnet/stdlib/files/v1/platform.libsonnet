{
  new(
    os='linux',
    arch='amd64',
    variant=null,
    kernel=null
  ):: {
    os: os,
    arch: arch,
    [if kernel != null then 'kernel']: kernel,
    [if variant != null then 'variant']: variant,
  },
}