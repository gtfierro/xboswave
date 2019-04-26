#!/bin/bash
set -ex

VERSION=$(cat ../../VERSION)

pushd ../../ingester
go build
popd
cp ../../ingester/ingester .
mkdir -p plugins
cp ../../ingester/plugins/*.so plugins/
docker build -t xbos/ingester:${VERSION} .
docker tag xbos/ingester:${VERSION} xbos/ingester:latest
docker push xbos/ingester:${VERSION}
docker push xbos/ingester:latest
