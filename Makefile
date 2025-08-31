.PHONY: clean all dev

DOCKER_COMPOSE_ARGS := --project-name part-tracker

all:
	docker compose $(DOCKER_COMPOSE_ARGS) up --detached --build

dev:
	docker compose $(DOCKER_COMPOSE_ARGS) watch

logs:
	docker compose $(DOCKER_COMPOSE_ARGS) logs -f

clean:
	docker compose $(DOCKER_COMPOSE_ARGS) down --volumes --remove-orphans
