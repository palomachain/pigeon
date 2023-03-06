#!/usr/bin/make -f

BRANCH               := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT               := $(shell git log -1 --format='%H')
TM_VERSION           := $(shell go list -m github.com/tendermint/tendermint | sed 's:.* ::')
BUILD_DIR            ?= $(CURDIR)/build
GOLANGCILINT_VERSION := 1.51.2

###############################################################################
##                                  Version                                  ##
###############################################################################

ifeq (,$(VERSION))
  VERSION := $(shell git describe --exact-match --tags 2>/dev/null | sed 's/^v//')
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=pigeon \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=pigeon \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/tendermint/tendermint/version.TMCoreSemVer=$(TM_VERSION)

ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

build: go.sum
	@echo "--> Building..."
	go build -mod=readonly $(BUILD_FLAGS) -o $(BUILD_DIR)/ ./...

build-linux: go.sum
	GOOS=linux GOARCH=amd64 $(MAKE) build

test:
	@echo "--> Testing..."
	@go test -v ./...

install-linter:
	@bash -c "source "scripts/golangci-lint.sh" && install_golangci_lint '$(GOLANGCILINT_VERSION)'"

lint: install-linter
	@echo "--> Linting..."
	@third_party/golangci-lint run --concurrency 16 ./...

.DEFAULT_GOAL := test
