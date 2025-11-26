.PHONY: migrate rollback fresh status seed seed-specific

# Migration commands
migrate:
	@go run cmd/migrate/main.go -command migrate

rollback:
	@go run cmd/migrate/main.go -command rollback -steps 1

rollback-steps:
	@go run cmd/migrate/main.go -command rollback -steps $(steps)

fresh:
	@go run cmd/migrate/main.go -command fresh

status:
	@go run cmd/migrate/main.go -command status

# Seeder commands
seed:
	@go run cmd/seed/main.go

seed-specific:
	@go run cmd/seed/main.go -seeder $(name)

# Help
help:
	@echo "Available commands:"
	@echo "  make migrate          - Run pending migrations"
	@echo "  make rollback         - Rollback last migration"
	@echo "  make rollback-steps=N - Rollback N migrations"
	@echo "  make fresh            - Drop all tables and re-run migrations"
	@echo "  make status           - Show migration status"
	@echo "  make seed             - Run all seeders"
	@echo "  make seed-specific name=seeder_name - Run specific seeder"

