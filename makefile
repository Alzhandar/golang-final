APP_NAME=restaurant
BUILD_DIR=./bin
MAIN_PATH=./cmd/app/main.go

.PHONY: build run clean test docker-up docker-down migrate-up migrate-down swagger lint

build:
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

run:
	@go run $(MAIN_PATH)

clean:
	@rm -rf $(BUILD_DIR)
	@go clean -cache

test:
	@go test ./... -v

docker-up:
	@docker-compose up -d

docker-down:
	@docker-compose down

migrate-up:
	@migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/restaurant?sslmode=disable" up

migrate-down:
	@migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/restaurant?sslmode=disable" down 1

swagger:
	@swag init -g $(MAIN_PATH) -o ./docs

lint:
	@golangci-lint run ./...
