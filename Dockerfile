FROM golang:1.23-alpine AS builder
WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Install gcc and musl-dev for sqlite3
ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev

# Build the server binary
RUN go build -v -o homelab_server homelab.com/homelab-server/homeLab-server/cmd

# SET UP sqlite db
COPY ./database ./database
RUN go run homelab.com/homelab-server/database/migration

FROM alpine:latest as final
WORKDIR /app

EXPOSE 6000
# TODO: IMPROVE env Injection
COPY .env .env

# Import database local.db
COPY --from=builder /app/database/local.db /app/database/local.db

# Create a group and user
RUN adduser -D myuser
USER myuser

# Run the binary
COPY --from=builder /app/homelab_server /app/homelab_server
CMD ["./homelab_server"]
