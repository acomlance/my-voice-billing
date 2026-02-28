# Reflection: Docker и удаление AuthKey

**Task ID:** docker-20250228  
**Date:** 2025-02-28

## Что сделано

- Реализован запуск проекта в Docker: Postgres 16 + приложение, docker-compose с healthcheck и initdb из `sql/docker-init` (01_accounts.sql, 02_tasks.sql). Конфиг для контейнеров — `config.docker.yml` с хостом `postgres`.
- Удалён AuthKey из домена: таблицы accounts (SQL), модель, репо, proto (RPC GetAccountByAuthKey убран), логика, хендлеры. Регенерация pb, сборка проверена.

## Уроки

- Инициализация схемы в Docker: порядок файлов в `docker-entrypoint-initdb.d` по алфавиту — 01_accounts, 02_tasks обеспечивает правильный порядок (tasks зависит от accounts).
- Конфиг в образе: копирование `config.docker.yml` в Dockerfile избавляет от монтирования в compose; для смены конфига без пересборки можно позже перейти на volume.

## Дальнейшее

- На Windows для запуска через Docker нужен Docker Desktop; без Docker — локальный Postgres и выполнение sql/public/*.sql.
- При изменении схемы accounts/tasks нужно синхронизировать sql/public и sql/docker-init.
