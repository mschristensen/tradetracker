# Trade Tracker

## Overview

Trade Tracker is a demo CLI application for ingesting trade data for a variety of instruments and tracking positions in those instruments.

### Features

Trade Tracker supports the following features:

- Plug into any stream of trade data. For demo purposes, a stream of random trades is used, but in a more realistic scenario the app could integrate with Kafka, a queue system etc.
- Trades are persisted to an append-only event store as timeseries data.
- Position data can be generated from those trades to understand how a position is changing with time.
- Query support to look up position size in an instrument at a given time.

### Trade Tracker Commands

- `tradetracker trade num instrumentID...` Simulates `num` random trades being streamed over a PubSub system.
- `tradetracker position intrumentID` (Re)generates position data from all trades for the given instrument.
- `tradetracker query intrumentID [timestamp]` Look up the position size at the given timestamp for an instrument. If no timestamp is provided, the latest position size is returned.

### Architecture

Trade Tracker consists of a CLI application backed by a PostgreSQL database for storing trades and positions.

Conceptually, Trade Tracker is composed of the following modules:

- The CLI tool entrypoint
- A `pubsub` module for simulating integration with a pub-sub system like Kafka.
- A `repo` module which provides an adapter for persisting trade and position data. This implementation uses PostgreSQL, but this could be swapped out e.g. a timeseries database.
- A `trade` module for consuming trade messages and writing them to the database via the repo.
- A `position` module for consuming trade messages, aggregating them to generate positions and writing them to the database via the repo.

### Project Structure

This project is organised in accordance with the best practices described [here](https://github.com/golang-standards/project-layout).

The folder structure is as follows:

- `/bin` holds any built executable binaries.
- `/cmd/<name>` holds the top-level entrypoints. Each folder `<name>` should hold a `main.go` to produce a standalone binary for that application.
- `/pkg` holds public packages that can be imported and used by external projects. Strict semantic versioning must be followed here.
- `/scripts` holds utility bash scripts and the like.
- `/configs` holds application config.
- `/internal` holds logic internal to the application which should not be imported by external projects (this is enforced by the Go compiler)
- `/internal/pkg` holds internal package code
- `/internal/app/apps` holds different application entrypoints. For example, a single binary can have multiple execution modes by leveraging different app implementations in here, while keeping dependencies separate.
- `/internal/app/cfg` holds application configuration to connect to external dependencies (databases, services, etc). Each configuration is implemented just once but can be used to configure any app type that needs it.

## Getting Started

### Database Setup

**Important Note:**

This has been developed and tested on MacOS Catalina (non-M1 chip). As long as you have `bash`, it should work but has not been tested.

First, you need PostgreSQL 12 running on port 5432. If you already have a native installation, you can use that. Otherwise, you can easily spin one up with Docker:

```
make pg_container
```

Once Postgres is running, you can setup the database with:

```
make database
```

You can then create the schema with:

```
make migrate direction=up
```

If you ever want to reset the database, you can run:

```
make migrate direction=down flags="-limit=0"
make migrate direction=up
```

### Application Setup

To build the application binary (you need Go 1.17):

```
make build
```

This puts the application binary in `./bin/tradetracker`. For ease of use, you can create an alias for your current session:

```
alias tradetracker="./bin/tradetracker"
```

Now all you need to do is set up your environment for the current shell session:

```
export $(cat ./configs/dev.env | xargs)
```

And now you're ready to start running `tradetracker` commands!

### Additional Help

To learn more about what you can do, run `make help`.

```
❯ make help
help:                          Shows help messages.
clean:                         Cleans up build artefacts.
database:                      Create the database and required roles.
migrate:                       Apply database migrations. [env, direction, flags]
lint:                          Runs linters.
docs:                          Starts the Go documentation server.
mocks:                         Generate mocks in all packages.
build_dependencies:            Builds the application dependencies.
build:                         Builds the application. [cmd]
reset_docker:                  Stops and cleans up running containers and volumes.
pg_container:                  Creates and runs a container running postgres.
test_integ_deps:               Prepares dependencies for running integration tests, creating and starting all containers.
test_start_containers:         Starts all the existing containers for the test environment.
test_integ:                    Runs integration tests. [timeout, dir, flags, run]
test_unit:                     Runs unit tests. [run, flags, timeout, dir]
exec:                          Executes the built application binary. [cmd, subcommand, flags]
run:                           Runs the application using go run. [cmd, subcommand, flags]
```

To learn more about what commands are available, run `tradetracker --help`:

```
❯ tradetracker --help
Usage:
   [flags]
   [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  position    Generates positions for an instrument from trade data after the given timestamp.
  query       Query for the position of an instrument at a given time.
  trade       Generates random trade data.

Flags:
      --env string                 Describes the current environment and should be one of: local, test, dev, prod. (default "local")
      --health_port int            The port the health server should listen on. (default 8080)
  -h, --help                       help for this command
      --log_level string           Sets the log level and should be one of: debug, info, warn, error. (default "debug")
      --max_goroutines int         The maximum allowed number of goroutines that can be spawned before healthchecks fail. (default 200)
      --max_pg_idle_conn int       The max number of allowed idle connections in the postgres connection pool. (default 80)
      --max_pg_open_conn int       The max number of allowed open connections in the postgres connection pool. (default 80)
      --port int                   The port the gRPC server should listen on. (default 8081)
      --postgres_database string   The database name for the postgres database. (default "tradetracker")
      --postgres_host string       The database host for the postgres database. (default "localhost")
      --postgres_password string   The database password for the postgres database. (default "tradetracker")
      --postgres_port int          The database port for the postgres database. (default 5432)
      --postgres_user string       The database user for the postgres database. (default "tradetracker")

Use " [command] --help" for more information about a command.
```

Or for help on a specific command:

```
❯ tradetracker trade --help
Generates random trade data.

Usage:
   trade num instrumentID... [flags]

Flags:
  -h, --help   help for trade

Global Flags:
      --env string                 Describes the current environment and should be one of: local, test, dev, prod. (default "local")
      --health_port int            The port the health server should listen on. (default 8080)
      --log_level string           Sets the log level and should be one of: debug, info, warn, error. (default "debug")
      --max_goroutines int         The maximum allowed number of goroutines that can be spawned before healthchecks fail. (default 200)
      --max_pg_idle_conn int       The max number of allowed idle connections in the postgres connection pool. (default 80)
      --max_pg_open_conn int       The max number of allowed open connections in the postgres connection pool. (default 80)
      --port int                   The port the gRPC server should listen on. (default 8081)
      --postgres_database string   The database name for the postgres database. (default "tradetracker")
      --postgres_host string       The database host for the postgres database. (default "localhost")
      --postgres_password string   The database password for the postgres database. (default "tradetracker")
      --postgres_port int          The database port for the postgres database. (default 5432)
      --postgres_user string       The database user for the postgres database. (default "tradetracker")
```

### Running Tests

**To run unit tests:**

```
make test_unit
```

**To run integration tests:**

Make sure everything is clean:

```
make reset_docker
```

Set up integration test dependencies:

```
make test_integ_deps
```

Run integration tests:

```
make test_integ
```


### Run a Demo

I have created an illustrative example of how to use `tradetracker`, which you can run with:

```
bash ./scripts/demo.sh
```

## Additional Notes

### Further Improvements

- Dockerise the CLI application. I didn't have time for this, but I hope you don't have too much trouble getting up and running.
- Comprehensive unit and integration testing. I didn't have time for this, but I included some small example tests as a demonstration.
- Implement a gRPC or HTTP server that the CLI tool can call to query for positions. This would allow the service to be deployed in the cloud and users can access the service from their own machines.
- Deploy the position and trade modules as stand alone services that can be scaled horizontally.
- Implement functionality to generating positions by aggregating trades over discrete time windows. This would allow a view of position data to be generated at the temporal granularity required by a given application. For example, for day traders who need high frequency updates, a small bin width could be used to generate positions with hgh granularity for a short time window. On the other hand, long term strategists could use a a larger bin width to generate positions with a lower granularity but spanning multiple years.
- Older trade data could be warehoused after long periods of time, according to business requirements.
  - As new trades messages come in, they could arrive out-of-order. However, at some point, all trades for a given time period will have been processed, allowing us to compute positions and freeze our view of the data.
- Use a timeseries database that supports partitioning. In the medium term, Postgres addons like Timescale could unlock some scale. Other technologies to explore include e.g. Amazon Timestream
