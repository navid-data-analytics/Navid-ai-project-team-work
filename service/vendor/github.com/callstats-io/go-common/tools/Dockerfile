FROM golang:1.8

WORKDIR /go/src/github.com/callstats-io/go-common/tools/
ADD vendor/github.com/callstats-io/go-common/tools/ .
ADD vendor vendor

RUN go install github.com/callstats-io/go-common/tools
