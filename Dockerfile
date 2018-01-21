FROM drone/ca-certs

ADD release/linux/amd64/drone /bin/

ENTRYPOINT ["/bin/drone"]
