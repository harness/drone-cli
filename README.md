# drone-cli

Command line client for the Drone continuous integration server.

Documentation: https://docs.drone.io/cli

Technical Support: https://discourse.drone.io

Bug Tracker: https://discourse.drone.io/c/bugs

## Release procedure

Run the changelog generator.

```BASH
docker run -it --rm -v "$(pwd)":/usr/local/src/your-app githubchangeloggenerator/github-changelog-generator -u drone -p drone-cli -t <secret github token>
```

You can generate a token by logging into your GitHub account and going to Settings -> Personal access tokens.

Next we tag the PR's with the fixes or enhancements labels. If the PR does not fufil the requirements, do not add a label.

** Before moving on make sure to update the version file `version/version.go`. **

Run the changelog generator again with the future version according to semver.

```BASH
docker run -it --rm -v "$(pwd)":/usr/local/src/your-app githubchangeloggenerator/github-changelog-generator -u drone -p drone-cli -t <secret token> --future-release v1.0.0
```

Create your pull request for the release. Get it merged then tag the release.
