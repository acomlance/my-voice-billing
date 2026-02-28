# Archive: Refactor gRPC errors, repo naming, config, logging

**Task ID:** refactor-20250228  
**Completed:** 2025-02-28

## Цель

Внедрить рекомендации по анализу проекта: исправить сборку (errors), уменьшить дублирование в хендлерах, унифицировать именование репо, вынести путь конфига в флаг/env, логировать Internal-ошибки.

## Выполненные изменения

### 1. gRPC errors (`internal/transport/grpc/errors/errors.go`)

- Экспорт `IsConflict` (ранее `isConflict`) для использования из handlers.
- Добавлена `ToStatus(err error) *status.Status`: NotFound → `codes.NotFound`, Conflict → `codes.AlreadyExists`, остальное → `codes.Internal`.

### 2. Handlers (`internal/transport/grpc/handlers/`)

- Введена общая `handleErr(err error, method string) error`: логирует ошибку при не NotFound и не Conflict, возвращает `errors.ToStatus(err).Err()`.
- `token.go`, `task.go`: все ветки с проверкой `IsNotFound`/`IsConflict` и `status.Error` заменены на `return nil, handleErr(err, "MethodName")`.
- В `CreateToken` исправлен возврат: `CreateTokenResponse` вместо `GetTokenResponse`.
- Удалены неиспользуемые импорты из `task.go` (errors, log).

### 3. Репозитории (`internal/repository/repo/`)

- `TokenService` → `TokenRepo`, `NewTokenService` → `NewTokenRepo`.
- `TaskService` → `TaskRepo`, `NewTaskService` → `NewTaskRepo`.
- Интерфейсы `TokenRepository`, `TaskRepository` без изменений.

### 4. Контейнер (`internal/container/container.go`)

- Поля `TokenRepo *repo.TokenService`, `TaskRepo *repo.TaskService` заменены на `*repo.TokenRepo`, `*repo.TaskRepo`.
- Вызовы `NewTokenService`/`NewTaskService` заменены на `NewTokenRepo`/`NewTaskRepo`.

### 5. Конфиг (`cmd/app/main.go`)

- Путь к конфигу: флаг `-config`, при отсутствии — `CONFIG_PATH`, по умолчанию `config/config.yml`.

## Проверка

- `go build -o build/app.exe ./cmd/app` — успешно.

## Связанные файлы

- `internal/transport/grpc/errors/errors.go`
- `internal/transport/grpc/handlers/token.go`
- `internal/transport/grpc/handlers/task.go`
- `internal/repository/repo/token_repo.go`
- `internal/repository/repo/task_repo.go`
- `internal/container/container.go`
- `cmd/app/main.go`
