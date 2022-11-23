FROM golang:1.19.3-alpine3.16
WORKDIR /app
COPY . /app 

RUN apk upgrade --update && \
    apk --no-cache add git && \
    apk add curl

# RUN go get -u github.com/cosmtrek/air && \
#     go build -o /go/bin/air github.com/cosmtrek/air
RUN go install github.com/cosmtrek/air@latest


RUN ls

CMD ["air", "-c", ".air.toml"]