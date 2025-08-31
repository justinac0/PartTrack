.PHONY: clean all dev

DOCKER_COMPOSE_ARGS := --project-name part-tracker

dev:
	docker compose $(DOCKER_COMPOSE_ARGS) watch

all:
	docker compose $(DOCKER_COMPOSE_ARGS) up --detached --build

logs:
	docker compose $(DOCKER_COMPOSE_ARGS) logs -f

clean:
	docker compose $(DOCKER_COMPOSE_ARGS) down --volumes --remove-orphans
