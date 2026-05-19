FROM golang:1.26-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/sync-server ./cmd/server

FROM alpine:3.22
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=build /out/sync-server /usr/local/bin/sync-server
COPY migrations ./migrations
ENV SYNC_SERVER_ADDR=:8080
ENV SYNC_SERVER_MIGRATIONS_DIR=/app/migrations
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 CMD wget -qO- http://127.0.0.1:8080/healthz >/dev/null || exit 1
CMD ["sync-server", "serve"]
