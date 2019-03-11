SERVER_NAME := veda-vscode-web-server
SHELL := /bin/bash
# VERSION := $(shell git describe --tags --abbrev=0)
# REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.name=$(NAME)' \
	-X 'main.version=$(VERSION)' \
	-X 'main.revision=$(REVISION)'

## Build client and server
build: build-client build-server

## Build WebExtension
build-client:
	npm run build

## Build server binary
build-server:
	go build -ldflags "$(LDFLAGS)" -o "build/VEDA for VSCode Web Server.app/Contents/MacOS/$(SERVER_NAME)" server/main.go
	go build -ldflags "$(LDFLAGS)" -o "build/VEDA for VSCode Web Server.app/Contents/MacOS/code-server-wrapper" server/code-server-wrapper.go

## Install dependencies
deps:

## Install tools required for development
deps-dev:
	go get golang.org/x/tools/cmd/goimports
	go get github.com/golang/lint/golint

## Run tests
test:
	go test

## Format source codes
fmt:
	goimports -l -w .
	go fmt ./...

## Lint
lint:
	go vet ./...
	golint ./...
	errcheck -ignoretests -blank ./...

.PHONY: build build-client build-server test fmt lint
