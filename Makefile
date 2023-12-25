# Ali Ahmadi 2023

include .envrc

## help: print this help message
.PHONY: help
help:
	@echo "Usage: "
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n "Are you sure? [y/n] " && read ans && [ $${ans:-N} = y ]

## run: run the cmd/api application
.PHONY: run
run:
	@go run ./cmd/api -dsn=${dsn}

## psql: connect to the database using psql
.PHONY: psql
psql:
	psql ${dsn}

## up: apply all up database migrations
.PHONY: up
up: confirm
	@echo "Running migrations..."
	@migrate -path=./migrations/ -database=${dsn} up

## migration name=$1: create a new database migration
.PHONY: migration
migration:
	@echo "creating migration files for ${name}..."
	migrate create -seq -ext=.sql -dir=./migrations/ ${name}

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo "formatting code..."
	go fmt ./...
	@echo "vetting code..."
	go vet ./...
	@echo "running tests..."
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo "tidying and verifying module dependencies..."
	go mod tidy
	go mod verify
	@echo "vendoring dependencies..."
	go mod vendor


# Build targets

## build: build the cmd/api application
.PHONY: build
build:
	@echo "Building cmd/api..."
	GOOS=windows GOARCH=amd64 go build -a -ldflags="-s" -o=./bin/windows_amd64/Meow ./cmd/api
	GOOS=linux GOARCH=amd64 go build -a -ldflags="-s" -o=./bin/linux_amd64/Meow ./cmd/api
	GOOS=linux GOARCH=mips64 go build -a -ldflags="-s" -o=./bin/linux_mips64/Meow ./cmd/api
	GOOS=linux GOARCH=arm64 go build -a -ldflags="-s" -o=./bin/linux_arm64/Meow ./cmd/api
