FROM golang:1.19-alpine3.16
ENV ROOT=/go/src/app
WORKDIR ${ROOT}

RUN apk update && \
    apk add git \
    bash \
    curl \
    build-base \
    alpine-sdk \
    musl-dev \
    libc-dev

RUN go install github.com/cosmtrek/air@latest 
RUN go install github.com/ramya-rao-a/go-outline@latest 
RUN go install github.com/nsf/gocode@latest
RUN go install github.com/acroca/go-symbols@latest
RUN go install github.com/fatih/gomodifytags@latest
RUN go install github.com/josharian/impl@latest
RUN go install github.com/haya14busa/goplay/cmd/goplay@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install golang.org/x/lint/golint@latest
RUN go install golang.org/x/tools/gopls@latest

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

RUN chmod +x ./wait-for-it.sh