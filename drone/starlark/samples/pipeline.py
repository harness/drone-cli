# cd drone/starlark/samples
# drone script --source pipeline.py --stdout

load('docker.py', 'docker')

def build(version):
  return {
    'name': 'build',
    'image': 'golang:%s' % version,
    'commands': [
      'go build',
      'go test',
    ]
  }

def main(ctx):
  if ctx['build']['message'].find('[skip build]'):
    return {
      'kind': 'pipeline',
      'name': 'publish_only',
      'steps': [
        docker('octocat/hello-world'),
      ],
    }
  return {
    'kind': 'pipeline',
    'name': 'build_and_publish',
    'steps': [
      build('1.11'),
      build('1.12'),
      docker('octocat/hello-world'),
    ],
  }
