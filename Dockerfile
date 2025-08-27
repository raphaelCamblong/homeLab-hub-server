FROM golang:1.25.0-alpine3.22 AS builder
WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

ARG TASKFILE_VERSION=latest

RUN go install github.com/go-task/task/v3/cmd/task@${TASKFILE_VERSION}

COPY . .

RUN task build:release

FROM alpine:latest as final
WORKDIR /app

EXPOSE 6000

COPY --from=builder /app/config.toml /app

# Create a group and user
RUN adduser -D myuser
USER myuser

# Run the binary
COPY --from=builder /app/homelab_server /app/homelab_server
CMD ["./homelab_server", "-c", " ./config.toml"]
