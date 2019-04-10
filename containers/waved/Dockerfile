FROM ubuntu:bionic

RUN apt-get update && apt-get install -y ca-certificates && \
    update-ca-certificates && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ENV GOTRACEBACK=all
ADD waved /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/waved", "--config", "/etc/waved/waved.toml"]
