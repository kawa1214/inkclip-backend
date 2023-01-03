FROM golang:1.19-alpine3.16 AS builder
ENV ROOT=/go/src/app
WORKDIR ${ROOT}

RUN apk update && apk add git && apk add bash && apk add make && apk add curl
COPY go.mod go.sum ./
RUN go mod download

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate /usr/bin/migrate

COPY . ${ROOT}
RUN {\
    echo "\
ENV=\
DB_DRIVER=\
DB_SOURCE=\
MIGRATION_URL=\
SERVER_ADDRESS=\
TOKEN_SECRET_KEY=\
ACCESS_TOKEN_DURATION=\
REFRESH_TOKEN_DURATION=\
MAIL_HOSTNAME=\
MAIL_PORT=\
MAIL_USERNAME=\
MAIL_PASSWORD=\
FRONT_URL=\
" > app.env;\
}

RUN CGO_ENABLED=0 GOOS=linux go build -o $ROOT/binary -buildvcs=false

EXPOSE 8080
CMD ["/go/src/app/binary"]