#FROM ubuntu:bionic
FROM golang:1.12.1-stretch

RUN apt-get update && apt-get install -y ca-certificates && \
    update-ca-certificates && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ENV GOTRACEBACK=all
ADD ingester /usr/local/bin/
ADD plugins /plugins
ENTRYPOINT ["/usr/local/bin/ingester","-config","/etc/ingester/ingester.yml"]
