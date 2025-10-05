BINARY?=scribe
.PHONY: build run test

build:
	go build -o bin/$(BINARY) ./cmd/scribe

run:
	go run ./cmd/scribe

test:
	go test ./...
