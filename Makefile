.SILENT:
.PHONY:
DOCKER_COMPOSE = docker-compose -f docker-compose.yml

install:
	$(DOCKER_COMPOSE) exec app go get -d ./...

build:
	$(DOCKER_COMPOSE) exec app go build -v ./cmd/oauth-proxy

run:
	$(DOCKER_COMPOSE) exec app /app/oauth-proxy

start: build run