{
  new(
    name,
    external=null,
    secretbox=null
  ):: {
    [name]: {
      [if external != null then 'external']: external,
      [if secretbox != null then 'secretbox']: secretbox,
    },
  },
}