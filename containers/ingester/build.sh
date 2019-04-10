#!/bin/bash
set -ex

VERSION=$(cat ../../VERSION)

pushd ../../ingester
go build
popd
cp ../../ingester/ingester .
mkdir -p plugins
cp ../../ingester/plugins/*.so plugins/
docker build -t ingester:${VERSION} .
docker tag ingester:${VERSION} ingester:latest
