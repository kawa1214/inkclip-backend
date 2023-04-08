[![codecov](https://codecov.io/gh/kawa1214/inkclip-backend/branch/main/graph/badge.svg?token=92PTEFPCPI)](https://codecov.io/gh/kawa1214/inkclip-backend)

## Setup

1. Docker

Start docker container.

```sh
docker-compose up -d
```

2. DevContainer

Select `Dev Containers: Reopen in Container` from the VSCode command palette.

## Command

- Create migrate file

```sh
migrate create -ext sql -dir db/migration xxx
```

- deploy

```sh
`fly deploy --local-only`
```

- Connection to DB

```sh
`flyctl proxy 5432 --app inkclip-backend-db`
```
