FROM golang:1.11

# configure installation variables
ARG SERVICE_VERSION=UNDEFINED

# Add required files, ignores test
WORKDIR /go/src/github.com/callstats-io/ai-decision/service/
ADD vendor/ vendor/
ADD scripts/docker/install_service.sh install_service.sh
ADD gen/ gen/
ADD migrations/ migrations/
ADD src/ src/

RUN bash ./install_service.sh $SERVICE_VERSION github.com/callstats-io/ai-decision/service/src/config

EXPOSE 13050 13051

CMD $GOPATH/bin/ai-decision-service
