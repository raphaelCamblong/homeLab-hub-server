FROM golang:1.24.6-alpine3.22 AS builder
WORKDIR /app

# Install dependencies

ARG APP_NAME=homelab_server
ARG MAIN_PATH=homelab.com/homelab-server/homeLab-server/cmd
ARG BUILD_ARGS="-ldflags='-w -s -extldflags=-static' -a -installsuffix cgo"

COPY . .

COPY go.mod go.sum ./
RUN go build ${BUILD_ARGS} -o ${APP_NAME} ${MAIN_PATH}

FROM alpine:3.22 as final
WORKDIR /app

EXPOSE 6000

COPY --from=builder /app/config.toml /app

# Create a group and user
RUN adduser -D myuser
USER myuser

# Run the binary
COPY --from=builder /app/homelab_server /app/homelab_server
CMD ["./homelab_server", "-c", " ./config.toml"]
