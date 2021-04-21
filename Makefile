# Go params
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=tensrose

DIST_FOLDER=dist
BIN_FOLDER=$(shell pwd)/bin

GOLANGCI_LINT_VERSION=1.39.0

.DEFAULT_GOAL := all

.PHONY: help
help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: all
all: build
build: ## build tensrose
	$(GOBUILD) -o $(BIN_FOLDER)/$(BINARY_NAME) -v cmd/tensrose/main.go

.PHONY: run-dev
run-dev: clean build ## run iris in development
	ulimit -n 1000
	./bin/reflex --decoration=fancy -r '\.go$$' -s -- sh -c 'make && $(BIN_FOLDER)/$(BINARY_NAME)'

.PHONY: clean
clean: ## clean package
	$(GOCLEAN)
	rm -rf $(DIST_FOLDER)

.PHONY: docker-build
docker-build: ## build docker images
	docker build -t iris:latest .

.PHONY: deploy
deploy: ## deploy to heroku
	@echo "Deploying to heroku"

.PHONY: build-all
build-all: clean ## build for all system and arch
	mkdir -p $(DIST_FOLDER)
	@echo "make a copy of dependencies"
	$(GOMOD) download
	# [darwin/amd64]
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(DIST_FOLDER)/$(BINARY_NAME)_darwin -v cmd/tensrose/main.go
	# [linux/amd64]
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(DIST_FOLDER)/$(BINARY_NAME)_linux -v cmd/tensrose/main.go
	# [windows/amd64]
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(DIST_FOLDER)/$(BINARY_NAME)_windows.exe -v cmd/tensrose/main.go

.PHONY: ensure-tools
ensure-tools: install-gofumports install-lint install-reflex ## ensure all dev tools

.PHONY: install-lint
install-lint:
	@echo "installing golangci-lint"
	if [ ! -x ./bin/golanci-lint ]  || ( ./bin/golangci-lint --version | grep -Fqv "version ${GOLANGCI_LINT_VERSION}" ); then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- v${GOLANGCI_LINT_VERSION}; \
	fi

.PHONY: install-gofumports
install-gofumports:
	if [ ! -x bin/gofumports ]; then \
		mkdir -p bin; \
		GOBIN=$(BIN_FOLDER) $(GOCMD) install mvdan.cc/gofumpt/gofumports@latest ; \
	fi

.PHONY: install-reflex
install-reflex:
	if [ ! -x bin/reflex ]; then \
		GOBIN=$(BIN_FOLDER) $(GOGET) github.com/cespare/reflex; \
	fi

.PHONY: lint
lint: install-lint ## lint
	./bin/golangci-lint run ./...

.PHONY: format
format: install-gofumports ## format
	find . -name \*.go | xargs ./bin/gofumports -local github.com/aarnphm/iris -w
