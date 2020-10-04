FROM alpine:3.12

LABEL maintainer="Joseph Salisbury <joseph@giantswarm.io>"

COPY ./prometheus-pingdom-exporter /prometheus-pingdom-exporter

RUN apk update && apk add ca-certificates

EXPOSE 8000

ENTRYPOINT ["/prometheus-pingdom-exporter"]
