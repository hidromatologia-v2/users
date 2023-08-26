# Users

[![Coverage](https://github.com/hidromatologia-v2/users/actions/workflows/codecov.yaml/badge.svg)](https://github.com/hidromatologia-v2/users/actions/workflows/codecov.yaml)
[![Release](https://github.com/hidromatologia-v2/users/actions/workflows/release.yaml/badge.svg)](https://github.com/hidromatologia-v2/users/actions/workflows/release.yaml)
[![Tagging](https://github.com/hidromatologia-v2/users/actions/workflows/tagging.yaml/badge.svg)](https://github.com/hidromatologia-v2/users/actions/workflows/tagging.yaml)
[![Test](https://github.com/hidromatologia-v2/users/actions/workflows/testing.yaml/badge.svg)](https://github.com/hidromatologia-v2/users/actions/workflows/testing.yaml)
[![codecov](https://codecov.io/gh/hidromatologia-v2/users/branch/main/graph/badge.svg?token=Q51HQV091I)](https://codecov.io/gh/hidromatologia-v2/users)
[![Go Report Card](https://goreportcard.com/badge/github.com/hidromatologia-v2/users)](https://goreportcard.com/report/github.com/hidromatologia-v2/users)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/hidromatologia-v2/users)

Users monolithic service API

## Documentation

| File                                                         | Description                                  |
| ------------------------------------------------------------ | -------------------------------------------- |
| [docs/spec.openapi.yaml](docs/spec.openapi.yaml)             | OpenAPI specification for this microservice. |
| [CONTRIBUTING.md](CONTRIBUTING.md)                           | Contribution manual.                         |
| [CICD.md](https://github.com/hidromatologia-v2/docs/blob/main/CICD.md) | CI/CD documentation.                         |

## Installation

### Docker

```shell
docker pull ghcr.io/hidromatologia-v2/users:latest
```

### Docker compose

```shell
docker compose -f ./docker-compose.dev.yaml up -d
```

### Binary

You can use the binary present in [Releases](https://github.com/hidromatologia-v2/users/releases/latest). Or compile your own with.

```shell
go install github.com/hidromatologia-v2/users@latest
```

## Config

| Variable              | Description                                                  | Example                                                      |
| --------------------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `MEMPHIS_STATION`     | Name for the station to **CREATE**/**CONNECT**               | `messages`                                                   |
| `MEMPHIS_PRODUCER`    | Alerts producer name                                         | `messages-producer`                                          |
| `MEMPHIS_HOST`        | Host or IP of the Memphis service                            | `10.10.10.10`                                                |
| `MEMPHIS_USERNAME`    | Memphis Username                                             | `root`                                                       |
| `MEMPHIS_PASSWORD`    | Memphis password, if ignored `MEMPHIS_CONN_TOKEN` will be used | `memphis`                                                    |
| `MEMPHIS_CONN_TOKEN`  | Memphis connection token, if ignored `MEMPHIS_PASSWORD` will be used | `ABCD`                                                       |
| `POSTGRES_DSN`        | Postgres DSN to be used                                      | `host=127.0.0.1 user=sulcud password=sulcud dbname=sulcud port=5432 sslmode=disable` |
| `JWT_SECRET`          | JWT secret to use                                            | `MY_SECRET`                                                  |
| `REDIS_EMAIL_ADDR`    | Address of the redis server                                  | `redis:6379`                                                 |
| `REDIS_EMAIL_DB`      | Redis database                                               | `1`                                                          |
| `REDIS_PASSWORD_ADDR` | Address of the redis server                                  | `redis:6379`                                                 |
| `REDIS_PASSWORD_DB`   | Redis database                                               | `2`                                                          |

### Binary

```shell
users HOST:PORT [HOST:PORT [...]]
```

## Coverage

| ![[coverage](https://app.codecov.io/gh/hidromatologia-v2/users)](https://codecov.io/gh/hidromatologia-v2/users/branch/main/graphs/sunburst.svg?token=Q51HQV091I) | ![[coverage](https://app.codecov.io/gh/hidromatologia-v2/users)](https://codecov.io/gh/hidromatologia-v2/users/branch/main/graphs/tree.svg?token=Q51HQV091I) |
| ------------------------------------------------------------ | ------------------------------------------------------------ |