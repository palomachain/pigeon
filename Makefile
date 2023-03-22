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

###############################################################################
##                            Helping target                                 ##
###############################################################################

RED_TEXT_BEGIN?=\033[0;31m
GREEN_TEXT_BEGIN?=\033[0;32m
YELLOW_TEXT_BEGIN?=\033[0;33m
AZURE_TEXT_BEGIN?=\033[0;36m
COLOURED_TEXT_END?=\033[0m

HELP_ROW_FORMAT := $(shell echo "$(AZURE_TEXT_BEGIN)%-30s$(COLOURED_TEXT_END)%s")

help::
	@echo "$(YELLOW_TEXT_BEGIN)Pigeon$(COLOURED_TEXT_END):"
	@printf "$(HELP_ROW_FORMAT)\n" "TARGET"         "DESCRIPTION"
	@echo "---------------------------------------------------------------------------------"
	@printf "$(HELP_ROW_FORMAT)\n" "build-linux"    "Builds the application binary for Linux AMD64"
	@printf "$(HELP_ROW_FORMAT)\n" "test"           "Runs go tests"
	@printf "$(HELP_ROW_FORMAT)\n" "lint"           "Lints (runs static checks) the code"


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
	@bash -c "source "scripts/golangci-lint.sh" && install_golangci_lint '$(GOLANGCILINT_VERSION)' '.'"

lint: install-linter
	@echo "--> Linting..."
	@third_party/golangci-lint run --concurrency 16 ./...

.DEFAULT_GOAL := help
