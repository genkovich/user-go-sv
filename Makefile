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

.PHONY: migration-up migration-down migration-status
