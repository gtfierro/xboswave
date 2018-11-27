GOPATH = /home/gabe/go
PLUGINS=$(wildcard plugins/*.go)

plugins/%.so: plugins/%.go FORCE
	go build -buildmode=plugin -o $@ $<

# this is a bit of a hack to make sure we always recompile the plugins
FORCE:

build: plugins/hamilton1.so plugins/iot_plugin.so
	go build

run: build 
	./ingester
