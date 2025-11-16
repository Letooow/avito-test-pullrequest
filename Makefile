SHELL = /bin/bash

.PHONY: build up down migrate logs ps restart

# Сборка образов
build:
	 docker compose build

# Явный прогон миграций (можно вызывать по отдельности)
migrate:
	 docker compose run --rm migrate

# Поднять всё: сначала миграции, затем сервисы
up: migrate
	docker compose up --build app

# Остановить и удалить контейнеры, сети и т.п.
down:
	docker compose down

# Перезапуск только приложения
restart:
	docker compose restart app

# Логи приложения
logs:
	docker compose logs -f app

# Статус сервисов
ps:
	docker compose ps
