# Account Service 

[![wercker status](https://app.wercker.com/status/9dad41bd24267b293467b812647f5d37/s/master "wercker status")](https://app.wercker.com/project/byKey/9dad41bd24267b293467b812647f5d37)


An account microservice that speaks gRPC backed by PostgreSQL, made with the [Lile generator](https://github.com/lileio/lile)

You can see the gRPC [protoc definition](https://github.com/arbarlow/account_service/blob/master/account/account.proto) for the RPC methods

The service will migrate and setup it's own tables if none exist in Postgres at the time of boot.

## ENVs

`DATABASE_URL` sets the PostgreSQL location.

```
DATABASE_URL="postgres://postgres@10.0.0.1/account_service"
```

## Docker

A pre build Docker container is available at:

```
docker pull lileio/account_service
```

## Development
The `docker-compose.yml` file will setup PostgreSQL with a default DB, but you will need create the test database

```
CREATE DATABASE account_service_test;
```
