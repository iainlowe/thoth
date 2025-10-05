BINARY?=daltu
.PHONY: build run test

build:
	go build -o bin/$(BINARY) ./cmd/daltu

run:
	go run ./cmd/daltu -addr=:8080

test:
	go test ./...
