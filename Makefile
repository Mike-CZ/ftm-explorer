# --------------------------------------------------------------------------
# Makefile for the Fantom Explorer
#
# (c) Fantom Foundation, 2023
# --------------------------------------------------------------------------

# project related vars
PROJECT := $(shell basename "$(PWD)")

# go related vars
GO_BASE := $(shell pwd)
GO_BIN := $(CURDIR)/build

# compile time variables will be injected into the app
APP_VERSION := 1.0
BUILD_DATE := $(shell date)
BUILD_COMPILER := $(shell go version)
BUILD_COMMIT := $(shell git show --format="%H" --no-patch)
BUILD_COMMIT_TIME := $(shell git show --format="%cD" --no-patch)

.PHONY: all clean test

all: bundle ftm-explorer

ftm-explorer:
	@go build \
    		-ldflags="-X 'ftm-explorer/cmd/ftm-explorer-cli/version.Version=$(APP_VERSION)' -X 'ftm-explorer/cmd/ftm-explore-cli/version.Time=$(BUILD_DATE)' -X 'ftm-explorer/cmd/ftm-explorer-cli/version.Compiler=$(BUILD_COMPILER)' -X 'ftm-explorer/cmd/ftm-explorer-cli/version.Commit=$(BUILD_COMMIT)' -X 'ftm-explorer/cmd/ftm-explorer-cli/version.CommitTime=$(BUILD_COMMIT_TIME)'" \
    		-o $(GO_BIN)/ftm-explorer \
    		-v \
    		./cmd/ftm-explorer-cli

bundle:
	@internal/api/graphql/schema/tools/make_bundle.sh

test:
	@go test ./...

clean:
	rm -fr ./build/*
