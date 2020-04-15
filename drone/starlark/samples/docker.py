def docker(repo):
  return {
    'name': 'docker',
    'image': 'plugins/docker',
    'settings': {
      'repo': repo,
    },
  }