GOPATH = /home/gabe/go
PLUGINS=$(wildcard plugins/*.go)

.PHONY: proto
proto: proto/xbos.proto
	protoc -Iproto/ -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:proto proto/*.proto
