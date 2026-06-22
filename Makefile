include ./.env

MIGRATION_PATH=db/migrations
DATABASE_URL=postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

migrate-create:
	@migrate create -ext sql -dir $(MIGRATION_PATH) -seq create_$(NAME)_table

migrate-up:
	@migrate -database $(DATABASE_URL) -path $(MIGRATION_PATH) up

migrate-down:
	@migrate -database $(DATABASE_URL) -path $(MIGRATION_PATH) down

migrate-force:
	@migrate -database "$(DATABASE_URL)" -path $(MIGRATION_PATH) force $(VERSION)

print-db-url:
	@echo "$(DATABASE_URL)"

seed:
	@psql "$(DATABASE_URL)" -f db/seeds/seed.sql