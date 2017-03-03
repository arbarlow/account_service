# Account Service

[![Build Status](https://travis-ci.org/lileio/account_service.svg?branch=master)](https://travis-ci.org/lileio/account_service) [![GoDoc](https://godoc.org/github.com/lileio/account_service?status.svg)](https://godoc.org/github.com/lileio/account_service)

An account microservice that speaks gRPC made with the [Lile generator](https://github.com/lileio/lile), backed by PostgreSQL or Cassandra.

``` protobuf
service AccountService {
  rpc List (ListAccountsRequest) returns (ListAccountsResponse) {}
  rpc GetById (GetByIdRequest) returns (Account) {}
  rpc GetByEmail (GetByEmailRequest) returns (Account) {}
  rpc AuthenticateByEmail (AuthenticateByEmailRequest) returns (Account) {}
  rpc Create (CreateAccountRequest) returns (Account) {}
  rpc Update (UpdateAccountRequest) returns (Account) {}
  rpc Delete (DeleteAccountRequest) returns (google.protobuf.Empty) {}
}
```
## Details

### Authentication
Passwords are stored hashed with [bcrypt](https://godoc.org/golang.org/x/crypto/bcrypt), no RPC method returns passwords or hashed passwords.

You can do simple authentication with the `AuthenticateByEmail` RPC method to roll your own authentication logic. I.e you can auth with email and password, but managing password length or auth tokens is up to you.

### Validations

At the moment the service will reject account create and update requests have either a blank name or email. "" is considered blank.

There is no email validation other than 'present' so to speak as I've never seen it done quite right.

## Docker

A pre build Docker container is available at:

```
docker pull lileio/account_service
```

## Commands

```
Usage:
  account_service [command]

Available Commands:
  migrate     Run database migrations
  server      Run the gRPC server
```

## Setup

Setup is configured via environment variables, depending on the database chosen.

The app provides migrations, but does create the databases itself.

### PostgreSQL

PostgreSQL is configured using the single ENV variable `POSTGRESQL_URL` and can be a url like string e.g.

`POSTGRESQL_URL="postgres://host/database"`

The PostgreSQL driver uses UUID's as primary key and a single table.

### Cassandra

Cassandra needs two ENV variables, the keyspace name and hosts to connect to (a comma seperated list):

`CASSANDRA_DB_NAME="account_service"`

`CASSANDRA_HOSTS="10.0.0.1,10.0.0.2"`

The Cassandra driver uses two tables so you lookup accounts by `ID` and by `Email`. See the migrations for more detailed information.


## Test
The `docker-compose.yml` file will run PostgreSQL and Cassandra for testing purposes, but you will need create the test databases yourself. Migrations are run automatically by the test suite.

For PostgreSQL (using psql):

``` sql
CREATE DATABASE account_service_test;
```

For Cassanda:

``` sql
create keyspace account_service WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
```

