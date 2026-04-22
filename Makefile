.PHONY: up down restart logs clean migrate-up migrate-down

# Load environment variables
include .env
export

# Spin up infrastructure with build and recreate
up:
	docker compose up -d --build --force-recreate

# Stop all services
down:
	docker compose down

# Restart services
restart: down up

# Follow backend logs
logs:
	docker compose logs -f backend

# Run database seeding script
seed:
	@chmod +x scripts/seed.sh
	@/bin/bash scripts/seed.sh
	
# Apply migrations (using 127.0.0.1 for host machine access)
migrate-up:
	migrate -path backend/migrations/ -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@127.0.0.1:5436/${POSTGRES_DB}?sslmode=disable" up

# Rollback all migrations
migrate-down:
	migrate -path backend/migrations/ -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@127.0.0.1:5436/${POSTGRES_DB}?sslmode=disable" down -all

# Force migration version
migrate-force:
	migrate -path backend/migrations/ -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@127.0.0.1:5436/${POSTGRES_DB}?sslmode=disable" force $(V)

# Full cleanup (removes volumes and local data)
clean:
	docker compose down -v
	rm -rf pgdata
