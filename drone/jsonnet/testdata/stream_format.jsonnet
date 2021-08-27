local pipeline(name) =
  {
    kind: 'pipeline',
    type: 'docker',
    name: name,
    platform: {
      os: 'linux',
      arch: 'amd64',
    },
    steps: [
      {
        name: 'test',
        image: 'golang:1.16',
        commands: ['go test ./...'],
      },
      {
        name: 'build',
        image: 'golang:1.16',
        commands: ['go build ./...'],
      },
    ],
  };

[
  pipeline('first'),
  pipeline('second'),
]
