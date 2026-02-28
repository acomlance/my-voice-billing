.PHONY: up down build start stop postgres postgres-reset

build:
	docker compose build app

postgres:
	docker compose up -d postgres

# Пересоздание контейнера и тома postgres — чистая БД, заново выполнятся sql/docker-init
postgres-reset:
	docker compose down -v
	docker compose up -d postgres

up: build
	docker compose up -d

down:
	docker compose down

start: up

stop: down
