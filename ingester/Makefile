PLUGINSRC=$(wildcard plugins/*.go)
PLUGINCOMP=$(PLUGINSRC:.go=.so)

# this is a bit of a hack to make sure we always recompile the plugins
FORCE: ;

plugins/%.so: plugins/%.go FORCE
	go build -buildmode=plugin -o $@ $<

build: $(PLUGINCOMP) FORCE
	go build

run: build 
	./ingester

prometheus:
	docker run --rm --name prometheus -p 9090:9090 -v `pwd`/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus

sshhostkeygen:
	ssh-keygen -f sshkey -N '' -t ed25519 -b 1024
