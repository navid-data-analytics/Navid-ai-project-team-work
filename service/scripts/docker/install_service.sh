SERVICE_BUILD_TIME="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"
FLAGS="-X $2.ServiceVersion=$1 -X $2.ImageBuildTime=$SERVICE_BUILD_TIME"
go build -ldflags "$FLAGS" -o $GOPATH/bin/ai-decision-service github.com/callstats-io/ai-decision/service/src
