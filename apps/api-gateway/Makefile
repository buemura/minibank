APP_NAME=api-gateway
MAIN_PATH=./cmd/server

.PHONY: deps start build test

install:
	go mod tidy

start:
	go run $(MAIN_PATH)

build:
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

test:
	go test ./...
