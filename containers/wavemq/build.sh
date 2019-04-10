#!/bin/bash
set -ex

VERSION=$(cat ../../VERSION)

pushd /home/$USER/go/src/github.com/immesys/wavemq
go build
popd
cp /home/$USER/go/src/github.com/immesys/wavemq/wavemq .
docker build -t wavemq:${VERSION} .
docker tag wavemq:${VERSION} wavemq:latest
