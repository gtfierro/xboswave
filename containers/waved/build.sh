#!/bin/bash
set -ex

VERSION=$(cat ../../VERSION)

pushd /home/$USER/go/src/github.com/immesys/wave/waved/cmd
go build
popd
cp /home/$USER/go/src/github.com/immesys/wave/waved/cmd/cmd waved
docker build -t waved:${VERSION} .
docker tag waved:${VERSION} waved:latest
