.PHONY: stop down start rebuild clean

# Остановить все контейнеры
stop:
	docker compose down

down: stop

# Удалить контейнеры и тома проекта, затем неиспользуемые контейнеры (если что-то повисло)
clean:
	docker compose down -v
	docker container prune -f

# 2. Запуск postgres + grpcui (app запускается из IDE с config.yml, БД localhost:5432)
start:
	docker compose up -d postgres grpcui

# Ребилд postgres (чистая БД, sql/docker-init) и подъём postgres+grpcui (http://localhost:9010)
rebuild:
	docker compose down -v
	docker compose up -d postgres grpcui
