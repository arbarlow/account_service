# Account Service 

[![wercker status](https://app.wercker.com/status/9dad41bd24267b293467b812647f5d37/s/master "wercker status")](https://app.wercker.com/project/byKey/9dad41bd24267b293467b812647f5d37) [![Go Report Card](https://goreportcard.com/badge/github.com/lileio/account_service)](https://goreportcard.com/report/github.com/lileio/account_service)

An account microservice that speaks gRPC made with the [Lile generator](https://github.com/lileio/lile), backed by PostgreSQL or Cassandra.

``` protobuf
service AccountService {
  rpc GetById (GetByIdRequest) returns (Account) {}
  rpc GetByEmail (GetByEmailRequest) returns (Account) {}
  rpc Create (CreateRequest) returns (Account) {}
  rpc Update (Account) returns (Account) {}
  rpc Delete (DeleteRequest) returns (DeleteResponse) {}
}
```

## Docker

A pre build Docker container is available at:

```
docker pull lileio/account_service
```

## Validations

At the moment the service will reject account create and update requests have either a blank name or email. "" is considered blank.

There is no email validation so to speak as I've never seen it done right.

## Setup

Setup is configured via environment variables, depending on the database chosen.

### PostgreSQL

PostgreSQL is configured using the single ENV variable `POSTGRESQL_URL` and can either be a url like string e.g. 

`POSTGRESQL_URL="postgres://host/database"`

 or a PostgreSQL connection string e.g:
 
 `POSTGRESQL_URL="host=localhost dbname=account_service sslmode=disable user=postgres"`
 
 Account service uses UUID's as primary key and a single table with the following schema:
 
 ``` sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS accounts (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
	name text NULL,
	email text NOT NULL,
	created_at timestamp without time zone NOT NULL DEFAULT (now() at time zone 'utc')
);
CREATE UNIQUE INDEX IF NOT EXISTS accounts_email ON accounts ((lower(email)));
 ```

### Cassandra

Cassandra needs two ENV variables, the keyspace name and hosts to connect to (a comma seperated list):

`CASSANDRA_DB_NAME="account_service"`

`CASSANDRA_HOSTS="10.0.0.1,10.0.0.2"`

Becuase of the way Cassandra works and uses primary keys, two tables are maintained so you lookup accounts by `ID` and by `Email`, it uses the follow schema:

``` sql
CREATE TABLE account_service.accounts_map_id (
    id text PRIMARY KEY,
    createdat timestamp,
    email text,
    name text
)

CREATE TABLE account_service.accounts_map_email (
    email text PRIMARY KEY,
    id text
)
```

## Development/Test
The `docker-compose.yml` file will run PostgreSQL and Cassandra, but you will need create the test databases

For PostgreSQL (using psql):

``` sql
CREATE DATABASE account_service_test;
```

For Cassanda:

``` sql
create keyspace account_service WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
```

