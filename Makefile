#!/usr/bin/make -f

GOLANGCILINT_VERSION := 1.51.2

test:
	@echo "--> Testing..."
	@go test -v ./...

install-linter:
	@bash -c "source "scripts/golangci-lint.sh" && install_golangci_lint '$(GOLANGCILINT_VERSION)'"

lint: install-linter
	@echo "--> Linting..."
	@third_party/golangci-lint run --concurrency 16 ./...

.DEFAULT_GOAL := test
