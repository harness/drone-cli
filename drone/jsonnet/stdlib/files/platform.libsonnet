{
  new(
    os='linux',
    arch='amd64',
    arm=null,
    version=null
  ):: {
    os: os,
    arch: arch,
    [if arm != null then 'arm']: arm,
    [if version != null then 'version']: version,
  },
}
