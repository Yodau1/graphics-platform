include .env
export

export PROJECT_ROOT=${shell pwd}

.DEFAULT_GOAL := help

env-up: ## Запустить PostgreSQL
	@docker compose up -d postgres
env-down:
	@docker compose down postgres
#read - чтение пользовательского ввода
# rm - удаление; -r - рекурсивное удаление; f - принудительное удаление, без подтверждения прав
#fi - конец ветвления
env-cleanup:
	@read -p "Очистить все volume файлы окружения? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down postgres && \
		rm -rf out/pgdata && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка отменена"; \
	fi;
env-port-forward: ## Открыть локальный порт PostgreSQL
	@docker compose up -d port-forwarder
env-port-close:
	@docker compose down port-forwarder
#run, а не up, потому что migrate отработает и сразу завершится,
#нет смысла на долгую запускать этот сервис
#--rm сразу удалить контейнер migrate
#"$(seq)" - динамическая передача переменной
#echo со статус-кодом 1 - есть ошибка; 0 - ошибок нет
migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует параметр seq. Пример: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm postgres-migrate \
	create \
	-ext sql \
	-dir /migrations \
	-seq "$(seq)"
migrate-up: ## Поднять миграции
	@make migrate-action action=up
migrate-down: ## Откатить все миграции
	@make migrate-action action=down
#make migrate-action action="down 1"
migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует параметр action. Пример: migrate-action action=up/down"; \
		exit 1; \
	fi; \
	docker compose run --rm postgres-migrate \
	-path /migrations \
	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable \
	"$(action)"
app-run: env-up env-port-forward ## Запустить приложение
	@ go mod tidy && \
	go run cmd/app/main.go

help: ## Показать доступные команды
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)