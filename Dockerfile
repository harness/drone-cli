FROM drone/ca-certs

COPY release/linux/amd64/drone /bin/

ENTRYPOINT ["/bin/drone"]
