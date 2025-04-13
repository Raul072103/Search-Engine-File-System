include .env
POSTGRES_DB_PATH = "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_IP}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"

.PHONY: migration
migration:
	@migrate create -seq -ext sql -dir ${POSTGRES_MIGRATIONS_PATH} $(word 2, $(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path $(POSTGRES_MIGRATIONS_PATH) -database $(POSTGRES_DB_PATH) up

.PHONY: migrate-down
migrate-down:
	@migrate -path $(POSTGRES_MIGRATIONS_PATH) -database $(POSTGRES_DB_PATH) down $(word 2, $(MAKECMDGOALS))

.PHONY: test-make-setup
test-make-setup:
	@echo ${POSTGRES_DB_PATH}