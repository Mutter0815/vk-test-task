# Makefile

# подхватываем .env (только переменные, без .env.example)
ifneq (,$(wildcard .env))
  include .env
  export
endif

COMPOSE = docker-compose

.PHONY: build up migrate logs down

# Собираем образы
build:
	$(COMPOSE) build --no-cache

# Запускаем БД и приложение (но не миграции)
up:
	$(COMPOSE) up -d postgres app

# Прогоняем миграции (up)
migrate:
	$(COMPOSE) run --rm migrate

# Смотреть логи всех сервисов
logs:
	$(COMPOSE) logs -f --tail=50

# Останавливаем и удаляем контейнеры и volume
down:
	$(COMPOSE) down -v
