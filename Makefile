PROJECT=prometheus-pingdom-exporter

BUILD_PATH := $(shell pwd)/.gobuild
GS_PATH := $(BUILD_PATH)/src/github.com/giantswarm
GOPATH := $(BUILD_PATH)

GOVERSION=1.6.2

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

clean:
	rm -rf $(BUILD_PATH) $(BIN) bin-dist/ build/

.gobuild:
	@mkdir -p $(GS_PATH)
	@rm -f $(GS_PATH)/$(PROJECT) && cd "$(GS_PATH)" && ln -s ../../../.. $(PROJECT)

	@builder get dep -b 1238ba19d24b0b9ceee2094e1cb31947d45c3e86 https://github.com/spf13/cobra.git $(GOPATH)/src/github.com/spf13/cobra
	@builder get dep -b cb88ea77998c3f024757528e3305022ab50b43be https://github.com/spf13/pflag.git $(GOPATH)/src/github.com/spf13/pflag
	
	@builder get dep -b cf63f55faae709a6bcd0ce28c4ae26f6106954cb https://github.com/russellcardullo/go-pingdom $(GOPATH)/src/github.com/russellcardullo/go-pingdom

	@builder get dep -b 488edd04dc224ba64c401747cd0a4b5f05dfb234 https://github.com/prometheus/client_golang.git $(GOPATH)/src/github.com/prometheus/client_golang
	@builder get dep -b 3ac7bf7a47d159a033b107610db8a1b6575507a4 https://github.com/beorn7/perks.git $(GOPATH)/src/github.com/beorn7/perks
	@builder get dep -b 3b06fc7a4cad73efce5fe6217ab6c33e7231ab4a https://github.com/golang/protobuf.git $(GOPATH)/src/github.com/golang/protobuf
	@builder get dep -b fa8ad6fec33561be4280a8f0514318c79d7f6cb6 https://github.com/prometheus/client_model.git $(GOPATH)/src/github.com/prometheus/client_model
	@builder get dep -b 3a184ff7dfd46b9091030bf2e56c71112b0ddb0e https://github.com/prometheus/common.git $(GOPATH)/src/github.com/prometheus/common
	@builder get dep -b abf152e5f3e97f2fafac028d2cc06c1feb87ffa5 https://github.com/prometheus/procfs.git $(GOPATH)/src/github.com/prometheus/procfs
	@builder get dep -b c12348ce28de40eed0136aa2b644d0ee0650e56c https://github.com/matttproud/golang_protobuf_extensions.git $(GOPATH)/src/github.com/matttproud/golang_protobuf_extensions

deps:
	@${MAKE} -B -s .gobuild

$(BIN): $(SOURCE) VERSION .gobuild
	CGO_ENABLED=0
	
	@echo Building inside Docker container for $(GOOS)/$(GOARCH)
	docker run \
	    --rm \
	    -v $(shell pwd):/usr/code \
	    -e GOPATH=/usr/code/.gobuild \
	    -e GOOS=$(GOOS) \
	    -e GOARCH=$(GOARCH) \
	    -w /usr/code \
	    golang:$(GOVERSION) \
	    $(BUILD_COMMAND)

ci-build: $(SOURCE) VERSION .gobuild
	CGO_ENABLED=0
	
	@echo Building for $(GOOS)/$(GOARCH)
	$(BUILD_COMMAND)
	
docker-image: $(BIN)
	docker build -t giantswarm/$(PROJECT):$(VERSION) .

bin-dist: $(SOURCE) VERSION .gobuild
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