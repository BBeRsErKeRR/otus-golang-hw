### Docker helper recieps
# ¯¯¯¯¯¯¯¯

DOCKER_MAKE_PATH:=$(abspath $(lastword $(MAKEFILE_LIST)))
DOCKER_MAKE_DIR:=$(dir $(DOCKER_MAKE_PATH))

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)
DOCKER_IMG="calendar:develop"


build-img:  ## Create docker image
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f $(DOCKER_MAKE_DIR)/build/Dockerfile .

run-img: build-img  ## Run  app container
	docker run $(DOCKER_IMG)

.PHONY: build-img run-img