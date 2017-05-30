# docker build --rm -t drone/drone .
FROM centurylink/ca-certs

ADD release/linux/amd64/drone /drone

ENTRYPOINT ["/drone"]
