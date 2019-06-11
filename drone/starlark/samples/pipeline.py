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
  if ctx['repo']['name'] == 'hello-world':
    print('hello-world')
  return {
    'kind': 'pipeline',
    'name': 'default',
    'steps': [
      build('1.11'),
      build('1.12'),
      docker('octocat/hello-world'),
    ],
  }
