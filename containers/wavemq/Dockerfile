FROM ubuntu:bionic

RUN apt-get update && apt-get install -y ca-certificates && \
    update-ca-certificates && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ENV GOTRACEBACK=all
ADD wavemq /usr/local/bin/
# mount /etc/wavemq
ENTRYPOINT ["/usr/local/bin/wavemq", "/etc/wavemq/wavemq.toml"]
