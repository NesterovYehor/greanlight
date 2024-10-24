# ==================================================================================== # 
# HELPERS
# ==================================================================================== #
#
## help: print this help message
.PHONY: help 
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'


.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ] 

# ==================================================================================== # 
# BUILD
# ==================================================================================== #


current_time = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/api ./cmd/api

# ==================================================================================== #  
# ENVIRONMENT LOADING
# ==================================================================================== #

# Load environment variables from .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# ==================================================================================== #  
# DEVELOPMENT
# ==================================================================================== #

# run/api: run the cmd/api application
.PHONY: run/api 
run/api:
	@go run ./cmd/api -db-dsn=${DATABASE_URL} -jwt-secret=${JWT_SECRET}


# db/psql: connect to the database using psql
.PHONY: db/psql 
db/psql:
	psql ${DATABASE_URL}


# db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm 
	migrate -path=./migrations -database="${DATABASE_URL}" up


# db/migrations/down name=$1: apply all down database migrations
.PHONY: db/migrations/down
db/migrations/down:
	migrate -path=./migrations -database="${DATABASE_URL}" down


# db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new: confirm
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext .sql -dir ./migrations ${name}

# ==================================================================================== # 
# QUALITY CONTROL
# ==================================================================================== #

# audit: tidy dependencies and format, vet and test all code
.PHONY: audit 
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

# vendor: tidy and vendor dependencies
.PHONY: vendor 
vendor:
	@echo 'Tidying and verifying module dependencies...' 
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor
