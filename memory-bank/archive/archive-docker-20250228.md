# Archive: Docker и удаление AuthKey

**Task ID:** docker-20250228  
**Completed:** 2025-02-28

## Цель

Поднять проект и Postgres в Docker; убрать AuthKey из сервиса.

---

## Выполненные изменения

### Docker

- **sql/docker-init/01_accounts.sql, 02_tasks.sql** — инициализация схемы при первом старте Postgres (без auth_key).
- **config/config.docker.yml** — хост БД `postgres`, порт 5432, user/postgres, password/root, database/voice_billing.
- **Dockerfile** — multi-stage: сборка в golang:1.26-alpine, запуск в alpine:3.20; в образ копируются бинарник и config.docker.yml; CMD с `-config /app/config/config.docker.yml`.
- **docker-compose.yml** — сервисы `postgres` (postgres:16-alpine, volume postgres_data, healthcheck pg_isready) и `app` (build: ., depends_on postgres condition: service_healthy, порт 50052:50052).
- **.dockerignore** — исключены .git, build/, logs/, memory-bank/ и др.

### Удаление AuthKey

- **SQL:** в sql/public/accounts.sql и sql/docker-init/01_accounts.sql удалены колонка auth_key и индекс accounts_auth_key_uindex.
- **models/account.go:** поле AuthKey удалено.
- **repo/account_repo.go:** метод GetByAuthKey и упоминания auth_key в запросах убраны; Create без auth_key.
- **proto:** RPC GetAccountByAuthKey и сообщение GetAccountByAuthKeyRequest удалены; из GetAccountResponse и CreateAccountRequest убрано поле auth_key. Выполнена регенерация pb.
- **logic/account_logic.go:** метод GetByAuthKey удалён.
- **handlers/account.go:** метод GetAccountByAuthKey удалён; CreateAccount и accountToProto без AuthKey.

## Проверка

- `go build -o build/app.exe ./cmd/app` — успешно.
- Запуск без Docker: локальный Postgres + sql/public/*.sql. Запуск с Docker: `docker compose up -d`.

## Связанные файлы

- docker-compose.yml, Dockerfile, .dockerignore
- config/config.docker.yml
- sql/docker-init/01_accounts.sql, 02_tasks.sql
- sql/public/accounts.sql
- internal/models/account.go
- internal/repository/repo/account_repo.go
- internal/services/logic/account_logic.go
- internal/transport/grpc/handlers/account.go
- proto/billing.proto, internal/transport/grpc/pb/*.go
