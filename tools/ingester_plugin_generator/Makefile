build:
	rm *plugin.go
	go build

.PHONY: generate
generate: build ingester_plugin_generator
	bash generate.sh
	cat hamilton_plugin_framework >> hamilton_plugin.go
	cat iot_plugin_framework >> iot_plugin.go
	gofmt -w *plugin.go
	cp *plugin.go ../../ingester/plugins/
