FROM drone/ca-certs
MAINTAINER Drone.IO Community <drone-dev@googlegroups.com>

LABEL org.label-schema.version=latest
LABEL org.label-schema.vcs-url="https://github.com/drone/drone-cli.git"
LABEL org.label-schema.name="Drone CLI"
LABEL org.label-schema.vendor="Drone.IO Community"
LABEL org.label-schema.schema-version="1.0"

ADD release/linux/amd64/drone /drone
ENTRYPOINT ["/drone"]