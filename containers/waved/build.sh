#!/bin/bash
set -ex

VERSION=$(cat ../../VERSION)

pushd /home/$USER/go/src/github.com/immesys/wave/waved/cmd
go build
popd
cp /home/$USER/go/src/github.com/immesys/wave/waved/cmd/cmd waved
docker build -t xbos/waved:${VERSION} .
docker tag xbos/waved:${VERSION} xbos/waved:latest
docker push xbos/waved:${VERSION}
docker push xbos/waved:latest
