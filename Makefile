.PHONY: down up rebuild clean

# Закрыть все контейнеры
down:
	docker compose down

# Запуск postgres + app (gRPC сервер) + grpcui
up:
	docker compose up -d postgres app grpcui

# Остановить контейнеры, пересобрать app без кэша, поднять стек
rebuild:
	make down
	docker compose build --no-cache app
	make up

# Удалить контейнеры и тома, неиспользуемые контейнеры (если что-то повисло)
clean:
	docker compose down -v
	docker container prune -f
