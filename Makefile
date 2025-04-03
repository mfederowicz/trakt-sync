.PHONY: test

export GO111MODULE=on

GIT_COMMIT ?= $(shell git rev-parse --verify HEAD)
GIT_VERSION ?= $(shell git describe --tags --always --dirty="-dev")
DATE ?= $(shell date -u '+%Y-%m-%d %H:%M UTC')
BUILDER ?= Makefile
VERSION_FLAGS := -X "trakt-sync/cli.date=$(DATE)" -X "trakt-sync/cli.builtBy=$(BUILDER)"

install:
	@go mod vendor

build:
	@go build -ldflags='$(VERSION_FLAGS)'

test:
	@go test -v -race ./...

cover:
	@go test -cover -coverprofile coverage.out ./...

linter:
	@revive --formatter friendly ./...

clean:
	@rm -rf *.json
