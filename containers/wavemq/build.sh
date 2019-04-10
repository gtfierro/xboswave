#!/bin/bash
set -ex

VERSION=$(cat ../../VERSION)

pushd /home/$USER/go/src/github.com/immesys/wavemq
go build
popd
cp /home/$USER/go/src/github.com/immesys/wavemq/wavemq .
docker build -t xbos/wavemq:${VERSION} .
docker tag xbos/wavemq:${VERSION} xbos/wavemq:latest
docker push xbos/wavemq:${VERSION}
docker push xbos/wavemq:latest
