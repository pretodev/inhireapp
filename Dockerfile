FROM golang:1.22.3 AS builder
WORKDIR  /app
COPY go.mod go.sum ./
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o inhire -a -ldflags '-linkmode external -extldflags "-static"' ./main.go

FROM alpine:3.20.2 AS binary
RUN apk update && apk add --no-cache chromium=128.0.6613.119-r0
WORKDIR /inhire
COPY --from=builder /app/inhire .
ENV SQLITE_DB_PATH=/inhire/db/inhire.sqlite
ENTRYPOINT ["./inhire"]
CMD ["version"]
