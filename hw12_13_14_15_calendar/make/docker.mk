### Docker helper recieps
# ¯¯¯¯¯¯¯¯

DOCKER_MAKE_PATH:=$(abspath $(lastword $(MAKEFILE_LIST)))
DOCKER_MAKE_DIR:=$(dir $(DOCKER_MAKE_PATH))

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)
DOCKER_IMG="calendar:develop"

define __COMPOSE_CMD
source $(DOCKER_MAKE_DIR).env && \
		docker-compose -f $(DOCKER_MAKE_DIR)docker-compose.yml
endef

.PHONY: run-img
build-img:  ## Create docker image
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f $(DOCKER_MAKE_DIR)/build/Dockerfile .

.PHONY: build-img
run-img: build-img  ## Run  app container
	docker run $(DOCKER_IMG)

.PHONY: service-up
service-up:
	@$(__COMPOSE_CMD) up -d $(_SERVICE)

.PHONY: service-stop
service-stop:
	@$(__COMPOSE_CMD) stop $(_SERVICE)

.PHONY: service-down
service-down:
	@$(__COMPOSE_CMD) down -v $(_SERVICE)

.PHONY: service-restart
service-restart:
	@$(__COMPOSE_CMD) restart $(_SERVICE)

.PHONY: service-attach
service-attach:
	@$(__COMPOSE_CMD) exec $(_SERVICE) sh

.PHONY: service-logs
service-logs:
	@$(__COMPOSE_CMD) logs --tail 100 $(_SERVICE)

.PHONY: service-status
service-status: ## See current services status
	@$(__COMPOSE_CMD) ps

.PHONY: calendar_postgres-up
calendar_postgres-up: ## Up dev postgress
	$(MAKE) _SERVICE=calendar_postgres service-up

.PHONY: calendar_postgres-stop
calendar_postgres-stop: ## Stop dev postgress
	$(MAKE) _SERVICE=calendar_postgres service-stop

.PHONY: calendar_postgres-logs
calendar_postgres-logs: ## See dev postgress logs
	$(MAKE) _SERVICE=calendar_postgres service-logs