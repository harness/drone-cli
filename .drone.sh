#!/bin/sh
set -e
set -x

# disable CGO for cross-compiling
export CGO_ENABLED=0

# compile for all architectures
GOOS=linux   GOARCH=amd64   go build -ldflags "-X main.version=${DRONE_TAG##v}" -o release/linux/amd64/drone       ./drone
GOOS=linux   GOARCH=arm64   go build -ldflags "-X main.version=${DRONE_TAG##v}" -o release/linux/arm64/drone       ./drone
GOOS=linux   GOARCH=ppc64le go build -ldflags "-X main.version=${DRONE_TAG##v}" -o release/linux/ppc64le/drone       ./drone
GOOS=linux   GOARCH=arm     go build -ldflags "-X main.version=${DRONE_TAG##v}" -o release/linux/arm/drone         ./drone
GOOS=windows GOARCH=amd64   go build -ldflags "-X main.version=${DRONE_TAG##v}" -o release/windows/amd64/drone.exe ./drone
GOOS=darwin  GOARCH=amd64   go build -ldflags "-X main.version=${DRONE_TAG##v}" -o release/darwin/amd64/drone      ./drone

# tar binary files prior to upload
tar -cvzf release/drone_linux_amd64.tar.gz   -C release/linux/amd64   drone
tar -cvzf release/drone_linux_arm64.tar.gz   -C release/linux/arm64   drone
tar -cvzf release/drone_linux_ppc64le.tar.gz -C release/linux/ppc64le drone
tar -cvzf release/drone_linux_arm.tar.gz     -C release/linux/arm     drone
tar -cvzf release/drone_windows_amd64.tar.gz -C release/windows/amd64 drone.exe
tar -cvzf release/drone_darwin_amd64.tar.gz  -C release/darwin/amd64  drone

# generate shas for tar files
sha256sum release/*.tar.gz > release/drone_checksums.txt
