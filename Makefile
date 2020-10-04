PROJECT=prometheus-pingdom-exporter

GOVERSION=1.15

BIN := $(PROJECT)

VERSION := $(shell cat VERSION)
COMMIT := $(shell git rev-parse --short HEAD)

.PHONY: all clean install

SOURCE=$(shell find . -name '*.go')

ifndef GOOS
  GOOS := $(shell go env GOOS)
endif
ifndef GOARCH
  GOARCH := $(shell go env GOARCH)
endif

BUILD_COMMAND=go build -a \
	-tags netgo \
	-ldflags \
	"-X github.com/giantswarm/prometheus-pingdom-exporter/cmd.version=$(VERSION) \
	-X github.com/giantswarm/prometheus-pingdom-exporter/cmd.goVersion=$(GOVERSION) \
	-X github.com/giantswarm/prometheus-pingdom-exporter/cmd.gitCommit=$(COMMIT) \
	-X github.com/giantswarm/prometheus-pingdom-exporter/cmd.osArch=$(GOOS)/$(GOARCH) \
	-w" \
	-o $(BIN)

all: $(BIN)

deps:
	GO111MODULE=on go get -v ./...

clean:
	rm -rf $(BIN) bin-dist/ build/

$(BIN): $(SOURCE) VERSION
	CGO_ENABLED=0
	
	@echo Building inside Docker container for $(GOOS)/$(GOARCH)
	docker run \
	    --rm \
	    -v $(shell pwd):/go/src/github.com/giantswarm/$(PROJECT) \
	    -e GOPATH=/go \
	    -e GOOS=$(GOOS) \
	    -e GOARCH=$(GOARCH) \
	    -e GO111MODULE=on \
	    -w /go/src/github.com/giantswarm/$(PROJECT) \
	    golang:$(GOVERSION) \
	    $(BUILD_COMMAND)

ci-build: $(SOURCE) VERSION
	CGO_ENABLED=0
	
	@echo Building for $(GOOS)/$(GOARCH)
	GO111MODULE=on $(BUILD_COMMAND)
	
docker-image: $(BIN)
	docker build -t giantswarm/$(PROJECT):$(VERSION) .

bin-dist: $(SOURCE) VERSION
	# Remove any old bin-dist or build directories
	rm -rf bin-dist build

	# Build for all supported OSs
	for OS in darwin linux; do \
		rm -f $(BIN); \
		GOOS=$$OS make $(BIN); \
		mkdir -p build/$$OS bin-dist; \
		cp README.md build/$$OS/; \
		cp LICENSE build/$$OS/; \
		cp $(BIN) build/$$OS/; \
		tar czf bin-dist/$(BIN).$(VERSION).$$OS.tar.gz -C build/$$OS .; \
	done

install: $(BIN)
	cp $(BIN) /usr/local/bin/