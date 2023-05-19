## simple makefile to log workflow
.PHONY: all test clean build install

SHELL := /bin/bash

GOFLAGS ?= $(GOFLAGS:)

all: dep install

build:
	@go build $(GOFLAGS) ./...

dep:
	@go get ./...

install:
	@go install $(GOFLAGS) ./...

clean:
	@go clean $(GOFLAGS) -i ./...

