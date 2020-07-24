.DEFAULT_GOAL := docker

# Docker-compose.yml path.
COMPOSE_PATH := ./deployment/docker-compose.yml

# Standarize go coding style for the whole app.
.PHONY: fmt
fmt:
	@go fmt ./...

# Clean project binary,
.PHONY: clean
clean:
	@go clean ./...

# Build docker images and containers.
.PHONY: docker-build
docker-build: clean fmt
	@docker-compose -f $(COMPOSE_PATH) -p point_system build
	@docker image prune -f --filter label=stage="ps_builder"

# Start built docker images and containers.
,PHONY: docker-up
docker-up:
	@docker-compose -f $(COMPOSE_PATH) -p point_system up -d
	@docker logs --follow point_system

# Build and start docker images and containers.
.PHONY: docker
docker: docker-build docker-up

# Tail docker worker logs.
.PHONY: docker-worker-logs
docker-worker-logs:
	@docker logs --follow ps_worker

# Stop all related docker containers.
.PHONY: docker-stop
docker-stop:
	@docker stop point_system ps_worker ps_db ps_rabbit

# Stop and remove all related docker container.
.PHONY: docker-rm
docker-rm:
	@docker rm point_system ps_worker ps_db ps_rabbit || echo ""
	@docker volume rm point_system_postgres-volume point_system_rabbitmq-volume || echo ""
	@docker rmi pointsystem ps_worker
