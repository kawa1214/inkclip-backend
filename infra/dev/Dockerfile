FROM golang:1.19-alpine3.16 AS builder
ENV ROOT=/go/src/app
WORKDIR ${ROOT}

RUN apk update && apk add git
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/cosmtrek/air@latest

COPY . ${ROOT}

EXPOSE 8080
CMD ["air", "-c", ".air.toml"]
