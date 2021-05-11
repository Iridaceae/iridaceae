# Go params
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# package-related
BINARY_NAME=iridaceae-server
TEST_BINARY_NAME=concertina-test
PKGDIR=cmd/iridaceae-server/main.go
TEST_PKGDIR=cmd/concertina-test/main.go 
PACKAGE_NAME=$(shell basename -s .git `git config --get remote.origin.url`)

# others
GOLANGCI_LINT_VERSION=1.39.0

# folders
DIST_FOLDER=dist
BIN_FOLDER=$(shell pwd)/bin

# Makefile settings
.DEFAULT_GOAL := all

.PHONY: help
help: ## display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[33m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: build ## run iridaceae in local context
	./bin/iridaceae-server

.PHONY: test
test:
	$(GOTEST) -v -race ./...

.PHONY: test-cov
test-cov:
	$(GOTEST) -v -race -covermode=atomic ./...

.PHONY: clean
clean:
	# creates /vendor
	$(GOMOD) tidy && $(GOMOD) vendor
	# then we clean
	$(GOCLEAN)
	rm -rf $(DIST_FOLDER)

.PHONY: dev
dev: clean ## run iris in development
	ulimit -n 1000
	./bin/reflex --decoration=fancy -r '\.go$$' -s -- sh -c 'make && $(BIN_FOLDER)/$(BINARY_NAME)'

.PHONY: all
all: build
build:
	$(GOBUILD) -o $(BIN_FOLDER)/$(BINARY_NAME) -v $(PKGDIR)
	$(GOBUILD) -o $(BIN_FOLDER)/$(TEST_BINARY_NAME) -v $(TEST_PKGDIR)

.PHONY: build-all
build-all: clean build docker-build ## build for all system and arch
	mkdir -p $(DIST_FOLDER)
	# [darwin/amd64]
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(DIST_FOLDER)/$(BINARY_NAME)_darwin -v $(PKGDIR)
	# [linux/amd64]
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(DIST_FOLDER)/$(BINARY_NAME)_linux -v $(PKGDIR)
	# [windows/amd64]
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(DIST_FOLDER)/$(BINARY_NAME)_windows.exe -v $(PKGDIR)

.PHONY: heroku-local
heroku-local: ## local deploy to heroku
	@echo "local deploy heroku"
	heroku local web

.PHONY: docker-build
docker-build: iris-build concertina-build ## build docker container

.PHONY: iris-build
iris-build:
	docker build --target=iridaceae-runner -t $(PACKAGE_NAME):latest .

.PHONY: concertina-build
concertina-build:
	docker build --target=concertina-runner -t concertina-go:latest .

.PHONY: docker-run
docker-run: iris-run concertina-run ## run docker container

.PHONY: iris-run
iris-run:
	docker run -t $(PACKAGE_NAME):latest

.PHONY: concertina-run
concertina-run:
	docker run -t concertina-go:latest

# TODO: Tags should just follow github revision or hash instead of latest.
.PHONY: docker-push
docker-push: iris-push concertina-push ## push container to docker registry

.PHONY: iris-push
iris-push: iris-build
	docker tag $(PACKAGE_NAME):latest aar0npham/iris-go:latest
	docker push aar0npham/iris-go:latest

.PHONY: concertina-push
concertina-push: concertina-build
	docker tag concertina-go:latest aar0npham/concertina-go:latest
	docker push aar0npham/concertina-go:latest

# TODO:
.PHONY: generate-env
generate-env: ## generate env file from defaults.example.env
	@./scripts/generate_env_file.sh

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
	find . -name \*.go | xargs ./bin/gofumports -local github.com/Iridaceae/iridaceae -w
	gofmt -w -s **/*.go
