.PHONY: up down restart logs migrate-up migrate-down clean

# Поднять всю инфраструктуру (БД, Бэкенд, Фронтенд)
up:
	docker compose up -d --build

# Остановить всё
down:
	docker compose down

# Перезапустить
restart: down up

# Посмотреть логи бэкенда
logs:
	docker compose logs -f backend

# Накатить миграции (запускать после 'make up')
# ВАЖНО: Убедись, что у друга установлен CLI утилита migrate, 
# либо он может запустить это локально, так как порт 5436 проброшен наружу
migrate-up:
	migrate -path migrations/ -database "postgres://postgres:secretpassword@127.0.0.1:5436/ancord_db?sslmode=disable" up

migrate-down:
	migrate -path migrations/ -database "postgres://postgres:secretpassword@127.0.0.1:5436/ancord_db?sslmode=disable" down -all

# Опасная зона: удалить контейнеры и базу данных (очистка)
clean:
	docker compose down -v
	rm -rf pgdata