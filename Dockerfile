FROM scratch
MAINTAINER Joseph Salisbury <joseph@giantswarm.io>

COPY ./prometheus-pingdom-exporter /prometheus-pingdom-exporter

EXPOSE 8000

ENTRYPOINT ["/prometheus-pingdom-exporter"]
