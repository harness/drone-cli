FROM alpine:3.7

RUN apk add --no-cache ca-certificates

COPY release/linux/amd64/drone /bin/

ENTRYPOINT ["/bin/drone"]
