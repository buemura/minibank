APP_NAME=account-service
MAIN_PATH=./cmd/server
MIGRATIONS_PATH=./migrations
DATABASE_URL?=postgres://postgres:postgres@localhost:5432/svc_account?sslmode=disable

.PHONY: deps start build test migrate-up migrate-down migrate-create

install:
	go mod tidy

start:
	go run $(MAIN_PATH)

build:
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

test:
	go test ./...

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $$name
