# Account Microservice [![wercker status](https://app.wercker.com/status/0f73245f410394e8b923cd22ca86970f/s/master "wercker status")](https://app.wercker.com/project/byKey/0f73245f410394e8b923cd22ca86970f)

An account microservice that speaks gRPC and is written in Go, backed by PostgreSQL.

You can see the gRPC [proto definition](https://github.com/arbarlow/account_service/blob/master/account/account.proto) for the RPC methods

Available on Docker.
```
docker pull alexrbarlow/account_service
```

## Development
The `docker-compose.yml` file will setup PostgreSQL with a default DB, but you will need create the test database

```
CREATE DATABASE account_service_test;
```

You can run the tests with Make
```
make test
```

## TODO
- Move to env variables
- Docker development env
- Prometheus
- Zipkin/Tracing
- Vendor dependencies
