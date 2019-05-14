PLUGINS=$(wildcard plugins/*.go)

.PHONY: proto
proto: proto/xbos.proto
	protoc -Iproto/ -Iproto/googleapis --go_out=plugins=grpc:proto proto/*.proto

.PHONY: proto-py
proto-py: wavemq/mqpb/wavemq.proto
	mkdir -p python/pyxbos/pyxbos/wavemq
	mkdir -p python/pyxbos/pyxbos/wave
	cp wavemq/mqpb/*.proto python/pyxbos/pyxbos/wavemq
	cp wave/eapi/pb/*.proto python/pyxbos/pyxbos/wave
	cd python/pyxbos; \
	poetry run python3 -m grpc_tools.protoc -Ipyxbos/wavemq -I../../proto/googleapis --python_out=pyxbos --grpc_python_out=pyxbos pyxbos/wavemq/*.proto; \
	poetry run python3 -m grpc_tools.protoc -Ipyxbos/wave -I../../proto/googleapis --python_out=pyxbos/wave --grpc_python_out=pyxbos/wave pyxbos/wave/*.proto; \
	poetry run python3 -m grpc_tools.protoc -I../../proto -I../../proto/googleapis --python_out=pyxbos --grpc_python_out=pyxbos ../../proto/*.proto; \
	sed -i -e 's/^import \(.*_pb2\)/from . import \1/g' pyxbos/*pb2*.py; \
	sed -i -e 's/^import \(.*_pb2\)/from . import \1/g' pyxbos/wave/*pb2*.py

#venv: python/requirements.txt
#	python3 -m venv venv; \
#	. venv/bin/activate; \
#	pip3 install -r python/requirements.txt
