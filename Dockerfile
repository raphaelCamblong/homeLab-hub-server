FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o homelab_server

FROM alpine:latest as final
WORKDIR /

COPY --from=builder /app/homelab_server /app/homelab_server

RUN adduser -D myuser
USER myuser

CMD ["./homelab_server"]
