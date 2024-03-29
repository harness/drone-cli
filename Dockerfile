FROM amd64/alpine:3.17 as alpine
RUN apk add -U --no-cache ca-certificates

FROM amd64/alpine:3.17
ENV GODEBUG netdns=go
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY release/linux/amd64/drone /bin/

ENTRYPOINT ["/bin/drone"]
