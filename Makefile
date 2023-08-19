# Makefile

include .env

migration-up:
	@GOOSE_DBSTRING=$(PSQL_DSN) GOOSE_DRIVER=postgres goose -dir ./migrations up

migration-down:
	@GOOSE_DBSTRING=$(PSQL_DSN) GOOSE_DRIVER=postgres  goose -dir ./migrations down

migration-status:
	@GOOSE_DBSTRING=$(PSQL_DSN)  GOOSE_DRIVER=postgres goose -dir ./migrations status

fixtures-up:
	@GOOSE_DBSTRING=$(PSQL_DSN) GOOSE_DRIVER=postgres goose -dir ./migrations/fixtures -table goose_fixtures up

fixtures-down:
	@GOOSE_DBSTRING=$(PSQL_DSN) GOOSE_DRIVER=postgres  goose -dir ./migrations/fixtures -table goose_fixtures down

docker-up:
	@cd docker && docker-compose up -d

docker-down:
	@cd docker && docker-compose down

setup-local-db: docker-up migration-up fixtures-up

.PHONY: migration-up migration-down migration-status fixtures-up fixtures-down docker-up
