#!/bin/sh

protoc -I=$GOPATH/src/github.com/callstats-io/ai-decision/service/protos \
    -I=$GOPATH/src/github.com/callstats-io/ai-decision/service/vendor \
    --go_out=plugins=grpc:$GOPATH/src/github.com/callstats-io/ai-decision/service/gen/protos/ \
    $GOPATH/src/github.com/callstats-io/ai-decision/service/protos/*.proto
