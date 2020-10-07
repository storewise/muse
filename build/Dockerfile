FROM golang:1 as builder
ENV GOOS=linux
ENV GOARCH=amd64
ARG WALRUS_VERSION=master

RUN git clone --branch ${WALRUS_VERSION} --depth=1 \
        https://github.com/lukechampine/walrus-cli.git "$GOPATH/src/lukechampine.com/walrus-cli"
WORKDIR $GOPATH/src/lukechampine.com/walrus-cli
RUN make


FROM debian:10-slim
LABEL maintainer="Junpei Kawamoto <junpei@storewise.tech>"

EXPOSE 9580
EXPOSE 9480
EXPOSE 9380
RUN apt-get update \
     && apt-get install -y --no-install-recommends curl \
     && apt-get -y clean \
     && rm -rf /var/lib/apt/lists/*
COPY bin/muse bin/musec build/readiness.sh /root/
COPY --from=builder /go/bin/* /root/
WORKDIR /root/

ENTRYPOINT ["/root/muse"]
CMD ["-d", "/data", "-serve-shard", "-s", "0.0.0.0:9480", "-serve-walrus", "-w", "0.0.0.0:9380"]
