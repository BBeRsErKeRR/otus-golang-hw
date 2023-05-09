### Docker helper recieps
# ¯¯¯¯¯¯¯¯

DOCKER_MAKE_PATH:=$(abspath $(lastword $(MAKEFILE_LIST)))
DOCKER_MAKE_DIR:=$(dir $(DOCKER_MAKE_PATH))

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)
DOCKER_IMG="calendar:develop"
DOCKER_BUILDKIT?=1
COMPOSE_DOCKER_CLI_BUILD?=1
DOCKER_REGISTRY?=
export

define __COMPOSE_CMD
source $(DOCKER_MAKE_DIR).env && \
		docker-compose -f $(DOCKER_MAKE_DIR)deployments/docker-compose.app.yml \
		-f $(DOCKER_MAKE_DIR)deployments/docker-compose.yml -p calendar
endef

bi:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		--build-arg=BIN_FILE="/opt/calendar/calendar-app" \
		--build-arg=MAIN_SRC="cmd/calendar/*" \
		--build-arg=CONFIG_SRC="configs/config.toml" \
		-t ${DOCKER_REGISTRY}calendar_app \
		-f ${DOCKER_MAKE_DIR}/build/app/Dockerfile ${DOCKER_CONTEXT_PATH}

.PHONY: run-img
build-img:  ## Create docker image
	docker-compose -f $(DOCKER_MAKE_DIR)deployments/docker-compose.app.yml build

.PHONY: build-img
run-img: build-img  ## Run  app container
	docker run $(DOCKER_IMG)

.PHONY: service-up
service-up:
	@$(__COMPOSE_CMD) up -d $(_SERVICE)

.PHONY: service-stop
service-stop:
	@$(__COMPOSE_CMD) stop $(_SERVICE)

.PHONY: service-restart
service-restart:
	@$(__COMPOSE_CMD) restart $(_SERVICE)

.PHONY: service-attach
service-attach:
	@$(__COMPOSE_CMD) exec $(_SERVICE) bash

.PHONY: service-logs
service-logs:
	@$(__COMPOSE_CMD) logs --tail 100 $(_SERVICE)

.PHONY: services-status
services-status: ## See current services status
	@$(__COMPOSE_CMD) ps

.PHONY: services-down
services-down: ## Remove all services
	@$(__COMPOSE_CMD) down -v

.PHONY: services-up
services-up: ## Up all services
	@$(__COMPOSE_CMD) up -d

.PHONY: postgres-up
postgres-up: ## Up dev postgress
	$(MAKE) _SERVICE=calendar_postgres service-up

.PHONY: postgres-stop
postgres-stop: ## Stop dev postgress
	$(MAKE) _SERVICE=calendar_postgres service-stop

.PHONY: postgres-logs
postgres-logs: ## See dev postgress logs
	$(MAKE) _SERVICE=calendar_postgres service-logs

.PHONY: postgres-attach
postgres-attach: ## Attach to psql container
	$(MAKE) _SERVICE=calendar_postgres service-attach

.PHONY: mq-up
mq-up: ## Up dev postgress
	$(MAKE) _SERVICE=calendar_mq service-up

.PHONY: postgres-stop
mq-stop: ## Stop dev postgress
	$(MAKE) _SERVICE=calendar_mq service-stop

.PHONY: postgres-logs
mq-logs: ## See dev postgress logs
	$(MAKE) _SERVICE=calendar_mq service-logs

.PHONY: postgres-attach
mq-attach: ## Attach to psql container
	$(MAKE) _SERVICE=calendar_mq service-attach


.PHONY: up
up:	build-img # Up all dev environment
	DURATION=$(or $(DURATION),120s) $(__COMPOSE_CMD) up -d

.PHONY: down # Down all dev environment
down: services-down