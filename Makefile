EXECUTABLES = git go basename pwd
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

#GOPATH=$(shell pwd)/vendor:$(shell pwd)
GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)
GONAME=$(shell basename "$(PWD)")
PLATFORMS=darwin linux
ARCHITECTURES=amd64

default: build

deps:
	@echo Ensuring Golang deps are up to date
	@dep ensure

build: deps
	@echo "Building $(GOFILES) to ./bin"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o bin/$(GONAME) $(GOFILES)

build_all: deps
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -v -o bin/$(GONAME)-$(GOOS)-$(GOARCH))))

run:
	@echo "Running $(GOFILES)"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go run $(GOFILES)

clean:
	@echo "Cleaning"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

.PHONY:	build build_all deps run clean