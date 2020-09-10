.POSIX:
.SUFFIXES:
.PHONY: all test clean

EXECUTABLE=git-topics
DISTDIR=dist
WINDOWS=$(DISTDIR)/$(EXECUTABLE)_windows_amd64.exe
LINUX=$(DISTDIR)/$(EXECUTABLE)_linux_amd64
DARWIN=$(DISTDIR)/$(EXECUTABLE)_darwin_amd64
VERSION=$(shell git describe --tags --always --long --dirty)

help: ## Show this help
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

all: test build # Build and run Tests

build: windows linux darwin ## Build binaries
	@echo version: $(VERSION)

test: ## Run the test suite
	go test -v ./...

clean: ## Remove previous build
	rm -f $(WINDOWS) $(LINUX) $(DARWIN)

install-dependencies: ## Install or update dependencies
	go mod tidy && go mod download

windows: $(WINDOWS) ## Build for Windows

linux: $(LINUX) ## Build for Linux

darwin: $(DARWIN) ## Build for Darwin (macOS)

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 go build -i -v -o $(WINDOWS) -ldflags="-s -w -X github.com/rody/find-commits/cmd/topics.version=$(VERSION)"  ./main.go

$(LINUX):
	env GOOS=linux GOARCH=amd64 go build -i -v -o $(LINUX) -ldflags="-s -w -X github.com/rody/find-commits/cmd/topics.version=$(VERSION)"  ./main.go

$(DARWIN):
	env GOOS=darwin GOARCH=amd64 go build -i -v -o $(DARWIN) -ldflags="-s -w -X github.com/rody/find-commits/cmd/topics.version=$(VERSION)"  ./main.go
