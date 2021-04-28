# Go params
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# others params
BINARY_NAME=tensrose
PACKAGE_NAME=$(shell basename -s .git `git config --get remote.origin.url`)
GOLANGCI_LINT_VERSION=1.39.0

# folders
DIST_FOLDER=dist
BIN_FOLDER=$(shell pwd)/bin

# Makefile settings
.DEFAULT_GOAL := all

.PHONY: help
help: ## display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[33m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: all
all: build
build:
	$(GOBUILD) -o $(BIN_FOLDER)/$(BINARY_NAME) -v pkg/cmd/tensroses-server/main.go

.PHONY: dev
dev: clean build ## run iris in development
	ulimit -n 1000
	./bin/reflex --decoration=fancy -r '\.go$$' -s -- sh -c 'make && make docker-build && $(BIN_FOLDER)/$(BINARY_NAME)'

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(DIST_FOLDER)

.PHONY: local-deploy
local-deploy: ## deploy to heroku
	@echo "local deploy heroku"
	heroku local web

.PHONY: docker-dev
docker-dev: docker-build docker-run  ## run development for ci

.PHONY: docker-build
docker-build:
	docker build -t $(PACKAGE_NAME):latest .

.PHONY: docker-run
docker-run:
	docker run -t $(PACKAGE_NAME):latest

.PHONY: build-all
build-all: clean build ## build for all system and arch
	mkdir -p $(DIST_FOLDER)
	# creates /vendor
	$(GOMOD) tidy && $(GOMOD) vendor
	# [darwin/amd64]
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(DIST_FOLDER)/$(BINARY_NAME)_darwin -v pkg/cmd/tensroses-server/main.go
	# [linux/amd64]
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(DIST_FOLDER)/$(BINARY_NAME)_linux -v pkg/cmd/tensroses-server/main.go
	# [windows/amd64]
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(DIST_FOLDER)/$(BINARY_NAME)_windows.exe -v pkg/cmd/tensroses-server/main.go

.PHONY: ensure-tools
ensure-tools: install-gofumports install-lint install-reflex ## ensure all dev tools

.PHONY: install-lint
install-lint:
	@echo "installing golangci-lint"
	if [[ ! -x bin/golangci-lint ]] || ( ./bin/golangci-lint --version | grep -Fqv "version ${GOLANGCI_LINT_VERSION}" ) ; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- v${GOLANGCI_LINT_VERSION}; \
	fi

.PHONY: install-gofumports
install-gofumports:
	if [[ ! -x bin/gofumports ]]; then \
		mkdir -p bin; \
		GOBIN=$(BIN_FOLDER) $(GOCMD) install mvdan.cc/gofumpt/gofumports@latest ; \
	fi

.PHONY: install-reflex
install-reflex:
	if [[ ! -x bin/reflex ]]; then \
		GOBIN=$(BIN_FOLDER) $(GOGET) github.com/cespare/reflex; \
	fi

.PHONY: ensure-format-lint
ensure-format-lint: format lint ## ensures you run format and lint

.PHONY: lint
lint: install-lint
	./bin/golangci-lint run ./...

.PHONY: format
format: install-gofumports
	find . -name \*.go | xargs ./bin/gofumports -local github.com/TensRoses/iris -w
