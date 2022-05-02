SHELL:=/bin/bash

cmd=tradetracker
subcommand=
postgres_image=postgres:12
timeout=30s
long_timeout=1m
dir=./...
run=.
flags=
direction=up
parallel=3


# include test env vars if the make target contains the string "test", otherwise use dev env vars
ifneq (true,$(NO_ENV_FILE))
ifneq (,$(findstring test,$(MAKECMDGOALS)))
include ./configs/test.env
export $(shell sed 's/=.*//' ./configs/test.env)
else ifeq (,$(findstring docker,$(MAKECMDGOALS))$(findstring docs,$(MAKECMDGOALS)))
include ./configs/dev.env
export $(shell sed 's/=.*//' ./configs/dev.env)
endif
endif

# set the postgres and rabbitmq images to be compatible with apple m1 chip
ifeq ($(shell uname -m),arm64)
postgres_image=arm64v8/postgres:12
endif


## help:				Shows help messages.
.PHONY: help
help:
	@sed -n 's/^##//p' $(MAKEFILE_LIST)


## clean:				Cleans up build artefacts.
.PHONY: clean
clean:
	@rm -rf bin
	@mkdir bin


## database:			Create the database and required roles.
.PHONY: database
database:
	@./scripts/db.sh


## migrate:			Apply database migrations. [env, direction, flags]
.PHONY: migrate
migrate:
	@go get github.com/rubenv/sql-migrate
	@go install github.com/rubenv/sql-migrate
	@sql-migrate $(direction) $(flags)
ifeq ($(direction),up)
	@./scripts/dump.sh
endif
	@go mod tidy


## lint:				Runs linters.
.PHONY: lint
lint:
	go fmt ./...
	go vet ./...
	golangci-lint run --timeout 5m0s ./...


## docs:				Starts the Go documentation server.
.PHONY: docs
docs:
	godoc -v -http=:6060


## mocks:				Generate mocks in all packages.
.PHONY:
mocks:
	@go install github.com/vektra/mockery/v2@latest
	@go generate ./...


## build_dependencies:		Builds the application dependencies.
.PHONY: build_dependencies
build_dependencies:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0
	@mkdir -p ./bin


## build:				Builds the application. [cmd]
.PHONY: build
build: build_dependencies lint
	@go build -a -installsuffix cgo -ldflags="-w -s" -o ./bin/$(cmd) ./cmd/$(cmd)/...
	@go mod tidy


## reset_docker:			Stops and cleans up running containers and volumes.
.PHONY: reset_docker
reset_docker:
	@-docker kill tradetracker_pg
	@-docker rm tradetracker_pg
	@-docker volume rm pg_tradetracker_data


## pg_container:			Creates and runs a container running postgres.
.PHONY: pg_container
pg_container:
	@echo "Spinning up postgres database..."
	@-docker volume create pg_tradetracker_data
	@docker pull $(postgres_image)
	@docker run -d --shm-size=2048MB -p $(POSTGRES_PORT):5432 --name tradetracker_pg -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -v pg_tradetracker_data:/var/lib/postgresql/data $(postgres_image)
	@echo "Waiting 60s for postgres..." && sleep 60
	@docker exec -it --user postgres tradetracker_pg psql -c "CREATE EXTENSION IF NOT EXISTS pgcrypto;"
	@-docker exec -it --user postgres tradetracker_pg psql -U postgres -c "CREATE USER $(POSTGRES_USER) WITH PASSWORD '$(POSTGRES_PASSWORD)' CREATEDB;"
	@-docker exec -it --user postgres tradetracker_pg psql -c "ALTER USER $(POSTGRES_USER) SUPERUSER;"
	@-docker exec -it --user postgres tradetracker_pg psql -c "CREATE DATABASE $(POSTGRES_DATABASE) OWNER $(POSTGRES_USER);"
	@echo "Postgres database ready"


## test_integ_deps:		Prepares dependencies for running integration tests, creating and starting all containers.
.PHONY: test_integ_deps
test_integ_deps: build_dependencies pg_container


## test_start_containers:		Starts all the existing containers for the test environment.
.PHONY: test_start_containers
test_start_containers:
	@docker start tradetracker_pg


## test_integ:			Runs integration tests. [timeout, dir, flags, run]
.PHONY: test_integ
test_integ: #test_start_containers
	@echo "Running integration tests on $(run)"
	@go test $(flags) -parallel $(parallel) -failfast -timeout $(long_timeout) $(dir) -run $(run)


## test_unit:			Runs unit tests. [run, flags, timeout, dir]
.PHONY: test_unit
test_unit: build_dependencies
	@echo "Running tests on $(run)"
	@go test $(flags) -short -parallel $(parallel) -failfast -timeout $(timeout) $(dir) -run $(run)


## exec:				Executes the built application binary. [cmd, subcommand, flags]
.PHONY: exec
exec:
	@./bin/$(cmd) $(subcommand) $(flags)


## run:				Runs the application using go run. [cmd, subcommand, flags]
.PHONY: run
run:
	@go run ./cmd/$(cmd) $(subcommand) $(flags)
