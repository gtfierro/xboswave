GOPATH = /home/gabe/go
PLUGINS=$(wildcard plugins/*.go)

.PHONY: proto
proto: proto/xbos.proto
	protoc -Iproto/ -Iproto/googleapis --go_out=plugins=grpc:proto proto/*.proto

.PHONY: wavemq-proto-py
wavemq-proto-py: venv wavemq/mqpb/wavemq.proto
	. venv/bin/activate; \
	python3 -m grpc_tools.protoc -Iwavemq/mqpb -Iproto -Iproto/googleapis --python_out=python/pyxbos --grpc_python_out=python/pyxbos wavemq/mqpb/*.proto
	
venv: python/requirements.txt
	python3 -m venv venv; \
	. venv/bin/activate; \
	pip3 install -r python/requirements.txt
