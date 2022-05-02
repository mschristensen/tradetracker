# Leather Wallet: Crypto.com Code Challenge

- [Leather Wallet: Crypto.com Code Challenge](#tradetracker-wallet-cryptocom-code-challenge)
  - [Overview](#overview)
    - [Features](#features)
    - [Leather Wallet Commands](#tradetracker-wallet-commands)
    - [Architecture](#architecture)
    - [How is a transaction executed?](#how-is-a-transaction-executed)
    - [Project Structure](#project-structure)
  - [Getting Started](#getting-started)
    - [Database Setup](#database-setup)
    - [Application Setup](#application-setup)
    - [Additional Help](#additional-help)
    - [Running Tests](#running-tests)
    - [Run a Demo](#run-a-demo)
  - [Additional Notes](#additional-notes)
    - [Further Improvements](#further-improvements)
    - [Time Spent](#time-spent)

## Overview

Leather Wallet is a demo CLI application offering users a multi-coin, online wallet and exchange experience.

Leather Wallet is also custodial. This approach offers a number of advantages:

1. It reduces the burden on the user as users are not responsible for their private and public keys.
2. The company operating Leather Wallet service can hold asset reserves in pooled wallets, and rebalance them on demand based on order book data.
3. Funds can be shifted between hot and cold wallets as dictated by security requirements.
4. Transactions can be batched to minimise on-chain transaction fees.
5. Cost efficiencies can be achieved through intelligent trading algorithms.

### Features

Leather Wallet supports the following features:

- Create a user.
- Create wallets for a given user. Each wallet holds a single currency.
- List the wallets, with their currencies and balances, for a user.
- Deposit money into a **fiat** wallet.
- Withdraw money from a **fiat** wallet.
- Transfer money in **any currency** from one wallet to another. Exchange functionality is provided to support transfers from one currency to another.

### Leather Wallet Commands

- `tradetracker create user [username]` Create a new user.
- `tradetracker create wallet [symbol] [username]` Create a wallet for a user (supported symbols: `GBP`, `BTC`, `ETH`).
- `tradetracker describe user [username]` List a user's wallets with their balances.
- `tradetracker describe wallet [wallet_id]` List the transactions associated with a wallet.
- `tradetracker deposit [amount] [wallet_id]` Deposit money in a wallet (`GBP` only).
- `tradetracker withdraw [amount] [wallet_id]` Withdraw money from a wallet (`GBP` only).
- `tradetracker send [amount] [from_wallet_id] [to_wallet_id]` Transfer funds from one wallet to another.

### Architecture

Leather Wallet consists of a CLI application backed by a PostgreSQL database for storing users, wallets and transactions. While Leather Wallet's database tracks the balances in the user wallets, this is in fact an abstraction over any real on-chain or other financial accounts.

Conceptually, Leather Wallet is composed of the following modules:

- The CLI tool entrypoint
- A `repo` module for interacting with the database
- A `broker` module for handling transaction requests
- An `exchange` module for executing transaction requests

### How is a transaction executed?

Say Alice wants to send 5 BTC to Bob's ETH wallet. The transaction is executed as follows:

1. The `broker` module checks whether Alice's has enough funds to complete this transaction. This considers not just her current wallet balance, but also its potential balance based on any pending (but not yet confirmed) transactions.
2. The `broker` module locks in an exchange rate for the user. It then creates two transactions in the *pending* state: the first one decrements Alice's wallet balance by 5 BTC, the second one increments Bob's wallet balance by the equivalent amount of ETH as determined by the exchange rate.
3. The `broker` then sends the transaction data to the `exchange`, which is responsible for making any necessary on-chain transactions, pooled wallet rebalancing, and the like. This asyncronous process concludes when the transaction is ultimately confirmed, and Alice and Bob's wallet balances are updated accordingly.


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
  create      Create a resource.
  deposit     Deposit money into a fiat wallet.
  describe    List information associated with a resource.
  help        Help about any command
  send        Send money from one wallet to another.
  withdraw    Withdraw money from a fiat wallet.

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
❯ tradetracker create --help
Create a resource.

Usage:
   create [flags]
   create [command]

Available Commands:
  user        Create a new user.
  wallet      Create a new wallet.

Flags:
  -h, --help   help for create

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

Use " create [command] --help" for more information about a command.
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

You can watch a video of this demo [here]().

## Additional Notes

### Further Improvements

- Dockerise the CLI application. I didn't have time for this, but I hope you don't have too much trouble getting up and running.
- Implement a gRPC or HTTP server that the CLI tool can call. This would allow the service to be deployed in the cloud and users can access the service from their own machines.
- Implement a user authentication system. Given the sensitive financial nature of the project, security would be a top priority and multi-factor authentication would be essential.
- Currently the `broker` simply synchronously sends transaction IDs to the `exchange`, which confirms them. In a more realistic scenario, the `broker` would asynchronously send more complete transaction data to the `exchange` over e.g. a queue. When the `exchange` has processed these transactions it would send a response message on a queue that would be picked up by a monitoring system. The monitoring system would then confirm or cancel transactions as needed. This has the additional advantage that transactions that occurred outside of the Leather Wallet system could be monitored and acted on accordingly.
- Currently the `broker` uses a fixed exchange rate provider, but it would be nice to interface with an API or on-chain Oracle to access live exchange rates.
- The `exchange` implementation is a dummy example for this project, but in reality could be responsible for:
  - Depositing or withdrawing fiat currency by interfacing with e.g. Open Banking APIs.
  - Sending crypto from one address to another
  - Batch processing on-chain transactions to obtain greater cost efficiency
  - Exchanging crypto using e.g. liquidity pools
  - Rebalancing pooled wallets
  - It should also be deployed as a standalone service.

### Time Spent

I quite enjoyed this code challenge, so ended up spending a bit longer on this than necessary. I completed this over two evenings, so I would say this took me about 10 hours.

