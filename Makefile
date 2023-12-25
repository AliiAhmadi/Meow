# Ali Ahmadi 2023

## help: print this help message
help:
	@echo "Usage: "
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

confirm:
	@echo -n "Are you sure? [y/n] " && read ans && [ $${ans:-N} = y ]

## run: run the cmd/api application
run:
	@go run ./cmd/api

## psql: connect to the database using psql
psql:
	psql ${dsn}

## up: apply all up database migrations
up: confirm
	@echo "Running migrations..."
	@migrate -path=./migrations/ -database=${dsn} up

## migration name=$1: create a new database migration
migration:
	@echo "creating migration files for ${name}..."
	migrate create -seq -ext=.sql -dir=./migrations/ ${name}
