### Calendar app automation recieps
# ¯¯¯¯¯¯¯¯

.DEFAULT_GOAL := help

SHELL=/bin/bash
SHELLOPTS:=$(if $(SHELLOPTS),$(SHELLOPTS):)pipefail:errexit

BIN := "./bin/calendar"


help: ## Display this help
	@IFS=$$'\n'; for line in `grep -h -E '^[a-zA-Z_#-]+:?.*?## .*$$' $(MAKEFILE_LIST)`; do if [ "$${line:0:2}" = "##" ]; then \
	echo $$line | awk 'BEGIN {FS = "## "}; {printf "\n\033[33m%s\033[0m\n", $$2}'; else \
	echo $$line | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'; fi; \
	done; unset IFS;

build: ## Build application
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build ## Run calendar app
	$(BIN) -config ./configs/config.toml

version: build  ## Project version
	$(BIN) version

test: ## Execute tests
	go test -race --tags=integration ./internal/...

coverage: ## test coverage
	go test --tags=integration -coverprofile=coverage.out ./...
	go tool cover -html coverage.out

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps ## Run liniter
	golangci-lint run --out-format=github-actions --path-prefix=hw12_13_14_15_calendar ./...

.PHONY: build run version test lint help

include make/docker.mk