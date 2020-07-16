# prometheus-pingdom-exporter

[![Build Status](https://api.travis-ci.org/giantswarm/prometheus-pingdom-exporter.svg)](https://travis-ci.org/giantswarm/prometheus-pingdom-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/giantswarm/prometheus-pingdom-exporter)](https://goreportcard.com/report/github.com/giantswarm/prometheus-pingdom-exporter)
[![GoDoc](https://godoc.org/github.com/giantswarm/prometheus-pingdom-exporter?status.svg)](http://godoc.org/github.com/giantswarm/prometheus-pingdom-exporter)
[![Docker](https://img.shields.io/docker/pulls/giantswarm/prometheus-pingdom-exporter.svg)](http://hub.docker.com/r/giantswarm/prometheus-pingdom-exporter) 

`prometheus-pingdom-exporter` exports Pingdom metrics to Prometheus.

## Prerequisites

## Getting `prometheus-pingdom-exporter`

Download the latest release: https://github.com/giantswarm/prometheus-pingdom-exporter/releases/latest

Clone the git repository: https://github.com/giantswarm/prometheus-pingdom-exporter.git

Download the latest docker image from here: https://hub.docker.com/r/giantswarm/prometheus-pingdom-exporter/


### How to build

#### Dependencies

- [github.com/russellcardullo/go-pingdom](https://github.com/russellcardullo/go-pingdom)
- [github.com/prometheus/client_golang](https://github.com/prometheus/client_golang)
- [github.com/spf13/cobra](https://github.com/spf13/cobra)

#### Building the binary

```
make
```

#### Building the docker image

```
make docker-image
```


## Running `prometheus-pingdom-exporter`

Running the binary directly:
```
$ prometheus-pingdom-exporter server <USERNAME> <PASSWORD> <API-KEY>
2016/07/09 12:20:51 Listening on port 8000
```

Running in a Docker container:
```
$ docker run -p 8000:8000 giantswarm/prometheus-pingdom-exporter:latest server <USERNAME> <PASSWORD> <API-KEY>
2016/07/09 11:22:04 Listening on port 8000
```

Help information can be found with the `--help` flag.

## Contact

- Mailing list: [giantswarm](https://groups.google.com/forum/!forum/giantswarm)
- Bugs: [issues](https://github.com/giantswarm/prometheus-pingdom-exporter/issues)

## Contributing & Reporting Bugs

See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches, the contribution workflow as well as reporting bugs.

## License

`prometheus-pingdom-exporter` is under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.
