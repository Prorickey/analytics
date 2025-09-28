FROM golang:1.25.1-alpine3.22 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download 

COPY . .

RUN go build -ldflags="-s -w" -o /usr/local/bin/app .

FROM alpine:3.22

WORKDIR /root/

COPY ./schema.sql /root/schema.sql

COPY --from=builder /usr/local/bin/app /usr/local/bin/app

EXPOSE 8080

CMD ["app"]