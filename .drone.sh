#!/bin/sh
set -e
set -x

mkdir bin
mkdir dist

# supported platforms
PLATFORMS=linux darwin windows

# compile drone for all architectures
GOOS=linux   GOARCH=amd64 go build -o ./bin/linux_amd64/drone   github.com/drone/drone-cli/drone
GOOS=darwin  GOARCH=amd64 go build -o ./bin/darwin_amd64/drone  github.com/drone/drone-cli/drone
GOOS=windows GOARCH=amd64 go build -o ./bin/windows_amd64/drone github.com/drone/drone-cli/drone

# tar binary files prior to upload
tar -cvzf dist/drone_linux_amd64.tar.gz   --directory=bin/linux_amd64   drone
tar -cvzf dist/drone_darwin_amd64.tar.gz  --directory=bin/darwin_amd64  drone
tar -cvzf dist/drone_windows_amd64.tar.gz --directory=bin/windows_amd64 drone

# generate shas for tar files
sha256sum ./dist/drone_linux_amd64.tar.gz   > ./dist/drone_linux_amd64.sha256
sha256sum ./dist/drone_darwin_amd64.tar.gz  > ./dist/drone_darwin_amd64.sha256
sha256sum ./dist/drone_windows_amd64.tar.gz > ./dist/drone_windows_amd64.sha256
