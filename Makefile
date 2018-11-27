GOPATH = /home/gabe/go
PLUGINS=$(wildcard plugins/*.go)

.PHONY: proto
proto: proto/xbos.proto
	protoc -Iproto/ -Iproto/googleapis --go_out=plugins=grpc:proto proto/*.proto

.PHONY: wavemq-proto-py
wavemq-proto-py: venv wavemq/mqpb/wavemq.proto
	mkdir -p python/pyxbos/wavemq
	mkdir -p python/pyxbos/wave
	cp wavemq/mqpb/*.proto python/pyxbos/wavemq
	cp wave/eapi/pb/*.proto python/pyxbos/wave
	. venv/bin/activate; \
	cd python/pyxbos; \
	python3 -m grpc_tools.protoc -Iwavemq -I../../proto/googleapis --python_out=. --grpc_python_out=. wavemq/*.proto; \
	python3 -m grpc_tools.protoc -Iwave -I../../proto/googleapis --python_out=wave --grpc_python_out=wave wave/*.proto; \
	python3 -m grpc_tools.protoc -I../../proto -I../../proto/googleapis --python_out=. ../../proto/*.proto; \
	sed -i -e 's/^import \(.*_pb2\)/from . import \1/g' *pb2*.py; \
	sed -i -e 's/^import \(.*_pb2\)/from . import \1/g' wave/*pb2*.py

venv: python/requirements.txt
	python3 -m venv venv; \
	. venv/bin/activate; \
	pip3 install -r python/requirements.txt
