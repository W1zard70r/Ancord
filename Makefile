.PHONY: up down restart logs clean migrate-up migrate-down

# Включаем чтение переменных из .env файла
include .env
export

# Поднять всю инфраструктуру (с пересборкой образов)
up:
	docker compose up -d --build --force-recreate

# Остановить всё
down:
	docker compose down

# Перезапустить
restart: down up

# Посмотреть логи бэкенда
logs:
	docker compose logs -f backend

seed:
	@chmod +x scripts/seed.sh
	@./scripts/seed.sh

# Накатить миграции (выполняется на хосте)
# Используем DB_URL из .env, но заменяем имя хоста 'postgres' на '127.0.0.1', 
# так как мигратор запускается локально (на твоем ПК), а не внутри сети Docker
migrate-up:
	migrate -path backend/migrations/ -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@127.0.0.1:5436/${POSTGRES_DB}?sslmode=disable" up

# Откатить миграции
migrate-down:
	migrate -path backend/migrations/ -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@127.0.0.1:5436/${POSTGRES_DB}?sslmode=disable" down -all

migrate-force:
	migrate -path backend/migrations/ -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@127.0.0.1:5436/${POSTGRES_DB}?sslmode=disable" force $(V)

# Полная очистка (удаляет БД!)
clean:
	docker compose down -v
	rm -rf pgdata