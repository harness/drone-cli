# drone-cli

Command line client for the Drone continuous integration server.

Documentation: https://docs.drone.io/cli

## Release procedure

Run the changelog generator.

```BASH
docker run -it --rm -v "$(pwd)":/usr/local/src/your-app githubchangeloggenerator/github-changelog-generator -u harness -p drone-cli -t <secret github token>
```

You can generate a token by logging into your GitHub account and going to Settings -> Personal access tokens.

Next we tag the PR's with the fixes or enhancements labels. If the PR does not fufil the requirements, do not add a label.

Run the changelog generator again with the future version according to semver.

```BASH
docker run -it --rm -v "$(pwd)":/usr/local/src/your-app githubchangeloggenerator/github-changelog-generator -u harness -p drone-cli -t <secret token> --future-release v1.0.0
```

Create your pull request for the release. Get it merged then tag the release.

## Community and Support

* [Harness Community](https://developer.harness.io/community)
* [Drone FAQ on Harness Developer Hub](https://developer.harness.io/kb/continuous-integration/drone-faqs).
* [Harness Community Slack](https://join.slack.com/t/harnesscommunity/shared_invite/zt-y4hdqh7p-RVuEQyIl5Hcx4Ck8VCvzBw) - Join the #drone slack channel to connect with our engineers and other users running Drone CI.
* [Drone on Meetup](https://www.meetup.com/harness/) - Check out previous events on [YouTube](https://www.youtube.com/watch?v=Oq34ImUGcHA&list=PLXsYHFsLmqf3zwelQDAKoVNmLeqcVsD9o).
