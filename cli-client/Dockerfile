FROM ubuntu:20.04

RUN apt-get update
RUN apt-get install -y wget git gcc

RUN wget -P /tmp https://go.dev/dl/go1.17.8.linux-amd64.tar.gz

RUN tar -C /usr/local -xzf /tmp/go1.17.8.linux-amd64.tar.gz
RUN rm /tmp/go1.17.8.linux-amd64.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

RUN mkdir -p /usr/src/casino-war/client

WORKDIR /usr/src/casino-war/client

# COPY . .

# RUN go build
