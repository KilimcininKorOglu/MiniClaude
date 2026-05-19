.PHONY: dev build test vet migrate-up migrate-down

BINARY := bin/sync-server

dev:
	go run ./cmd/server serve

build:
	mkdir -p bin
	go build -o $(BINARY) ./cmd/server

test:
	go test ./...

vet:
	go vet ./...

migrate-up:
	go run ./cmd/server migrate up

migrate-down:
	go run ./cmd/server migrate down
