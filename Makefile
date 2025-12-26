.PHONY: all clean clean-coverage generate install install-dependencies install-tools lint test test-verbose test-verbose-with-coverage

export ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
export PKG := github.com/avito-tech/go-mutesting
export ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
export BIN_NAME := go_mutesting
export OS_NAME := $(shell uname | tr A-Z a-z)
export COMMIT := $(shell git rev-parse HEAD)
export SECRET := 'secret'
export UPLOAD_TO := 'upload_url'

export TEST_TIMEOUT_IN_SECONDS := 240

$(eval $(ARGS):;@:) # turn arguments into do-nothing targets
export ARGS

ifdef ARGS
	PKG_TEST := $(ARGS)
else
	PKG_TEST := $(PKG)/...
endif

all: install-dependencies install-tools install lint test
.PHONY: all

clean:
	go clean -i $(PKG)/...
	go clean -i -race $(PKG)/...
.PHONY: clean

clean-coverage:
	find $(ROOT_DIR) | grep .coverprofile | xargs rm
.PHONY: clean-coverage

generate: clean
	go generate $(PKG)/...
.PHONY: generate

install:
	go install -v $(PKG)/...
.PHONY: install

install-dependencies:
	go get -t -v $(PKG)/...
.PHONY: install-dependencies

install-tools:
	# generation
	go install golang.org/x/tools/cmd/stringer

	# linting
	go install golang.org/x/lint/golint@latest
	go install github.com/kisielk/errcheck@latest

	# code coverage
	go install github.com/onsi/ginkgo/ginkgo@latest
	go install github.com/modocache/gover@latest
	go install github.com/mattn/goveralls@latest
.PHONY: install-tools

lint: ci-errcheck ci-gofmt ci-govet ci-lint
.PHONY: lint

test:
	go test -race -test.timeout "$(TEST_TIMEOUT_IN_SECONDS)s" $(PKG_TEST)
.PHONY: test

test-verbose:
	go test -race -test.timeout "$(TEST_TIMEOUT_IN_SECONDS)s" -v $(PKG_TEST)
.PHONY: test-verbose

test-verbose-with-coverage:
	ginkgo -r -v --cover --coverprofile=coverage.out --covermode=atomic --race -skipPackage="testdata"
.PHONY: test-verbose-with-coverage

ci-errcheck:
	$(ROOT_DIR)/scripts/ci/errcheck.sh
.PHONY: ci-errcheck

ci-gofmt:
	$(ROOT_DIR)/scripts/ci/gofmt.sh
.PHONY: ci-gofmt

ci-govet:
	$(ROOT_DIR)/scripts/ci/govet.sh
.PHONY: ci-govet

ci-lint:
	$(ROOT_DIR)/scripts/ci/lint.sh
.PHONY: ci-lint

.PHONY: build-all
build-all:
	make build GOOS=linux GOARCH=amd64
	make build GOOS=darwin GOARCH=amd64
	make build GOOS=darwin GOARCH=arm64

.PHONY: build
build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -buildvcs=false -ldflags="-s -w -X 'main.commit=$(COMMIT)'" $(TAGS) -o ./build/$(BIN_NAME)_$(GOOS)_$(GOARCH) ./cmd/go-mutesting

.PHONY: upload-all
upload-all:
	make upload GOOS=linux GOARCH=amd64
	make upload GOOS=darwin GOARCH=amd64
	make upload GOOS=darwin GOARCH=arm64

.PHONY: upload
upload:
	curl -u $(SECRET) $(UPLOAD_TO) -T "./build/$(BIN_NAME)_$(GOOS)_$(GOARCH)"