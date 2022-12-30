FROM golang:1.19-alpine3.16
ENV ROOT=/go/src/app
WORKDIR ${ROOT}

RUN apk update && apk add git && apk add bash
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/cosmtrek/air@latest

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

RUN chmod +x ./wait-for-it.sh