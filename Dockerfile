FROM golang:1.19-alpine3.16 AS builder
ENV ROOT=/go/src/app
WORKDIR ${ROOT}

RUN apk update && apk add git
COPY go.mod go.sum ./
RUN go mod download

COPY . ${ROOT}
# RUN touch app.env
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
" > app.env;\
}

RUN CGO_ENABLED=0 GOOS=linux go build -o $ROOT/binary -buildvcs=false

EXPOSE 8080
CMD ["/go/src/app/binary"]