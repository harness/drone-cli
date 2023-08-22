FROM --platform=$BUILDPLATFORM alpine:3.18 as alpine

RUN apk add -U --no-cache ca-certificates

FROM scratch

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ARG TARGETPLATFORM

COPY release/${TARGETPLATFORM}/drone /bin/

ENTRYPOINT ["/bin/drone"]
