### Calendar app automation recieps
# ¯¯¯¯¯¯¯¯

.DEFAULT_GOAL := help

BIN := "./bin/calendar"

help: SHELL=/bin/bash
help: ## Display this help
	@IFS=$$'\n'; for line in `grep -h -E '^[a-zA-Z_#-]+:?.*?## .*$$' $(MAKEFILE_LIST)`; do if [ "$${line:0:2}" = "##" ]; then \
	echo $$line | awk 'BEGIN {FS = "## "}; {printf "\n\033[33m%s\033[0m\n", $$2}'; else \
	echo $$line | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'; fi; \
	done; unset IFS;

generate: ## Generate proto files
	go $@ ./...

build: ## Build application
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build ## Run calendar app
	source make/.env && $(BIN) --config ./configs/config.toml

version: build  ## Project version
	$(BIN) version

test: ## Execute tests
	go test -race ./internal/... ./api/...

wait-for-db:
	sleep 5

integration: postgres-up wait-for-db migrate ## Execute integration tests
	go test --tags=$@ -race ./internal/...


coverage: postgres-up wait-for-db migrate ## test coverage
	go test --tags=integration -coverprofile=coverage.out ./internal/...
	go tool cover -html coverage.out

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps ## Run liniter
	golangci-lint run --config=$$(pwd)/../.golangci.yml \
		--timeout 3m0s \
		--skip-dirs='/opt/hostedtoolcache/go|/go/pkg/mod' \
		--path-prefix=hw12_13_14_15_calendar ./internal/... ./api/... ./cmd/...

goose:
	@(goose --version) || go install github.com/pressly/goose/v3/cmd/goose@latest

define __GOOSE_CMD
source make/.env && \
	goose --dir migrations postgres "postgres://$${DB_USERNAME}:$${DB_PASSWORD}@localhost:5532/$${DB_DATABASE}?sslmode=disable"
endef

migrate: ## Update db items 
	$(__GOOSE_CMD) up

downgrade: ## Downgrade db items
	$(__GOOSE_CMD) down

migration_status: ## Current migrations status
	$(__GOOSE_CMD) status

.PHONY: build run version test lint help coverage goose

include make/docker.mk