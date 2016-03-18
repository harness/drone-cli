# Docker image for the Drone build runner
# Refer to README.md for instructions on how to build the image

FROM gliderlabs/alpine:3.1
RUN apk-install ca-certificates && rm -rf /var/cache/apk/*
ADD drone-exec /bin/
ENTRYPOINT ["/bin/drone-exec"]
