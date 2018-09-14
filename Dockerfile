FROM alpine:3.8
MAINTAINER Joseph Salisbury <joseph@giantswarm.io>

COPY ./prometheus-pingdom-exporter /prometheus-pingdom-exporter

RUN apk update && apk add ca-certificates

EXPOSE 8000

ENTRYPOINT ["/prometheus-pingdom-exporter"]
