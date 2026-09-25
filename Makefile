.PHONY: test

export GO111MODULE=on

GIT_COMMIT ?= $(shell git rev-parse --verify HEAD)
GIT_VERSION ?= $(shell git describe --tags --always --dirty="-dev")
DATE ?= $(shell date -u '+%Y-%m-%d %H:%M UTC')
BUILDER ?= Makefile
VERSION_FLAGS := -X "github.com/mfederowicz/trakt-sync/cli.date=$(DATE)" -X "github.com/mfederowicz/trakt-sync/cli.builtBy=$(BUILDER)" -X "github.com/mfederowicz/trakt-sync/cli.version=$(GIT_VERSION)" -X "github.com/mfederowicz/trakt-sync/cli.commit=$(GIT_COMMIT)"

install:
	@go mod vendor

build:
	@go build -ldflags='$(VERSION_FLAGS)'

test:
	@go test -v -race ./...

cover:
	@go test -cover -coverprofile coverage.out ./...

linter:
	@revive --config ./revive.toml --formatter friendly ./...

cleanup:
	@find . -type f -name '*.go' -exec gofmt -w {} +

clean:
	@rm -rf *.json

release-notes:
	@awk -v ver="$(VERSION)" -f .github/release-notes.awk CHANGELOG.md
