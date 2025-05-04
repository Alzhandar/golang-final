APP_NAME=restaurant
BUILD_DIR=./bin
MAIN_PATH=./cmd/app/main.go

.PHONY: build run clean test docker-up docker-down migrate-up migrate-down swagger lint docker-migrate-up docker-migrate-down

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
	@migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/restaurant_db?sslmode=disable" up

migrate-down:
	@migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/restaurant_db?sslmode=disable" down 1

docker-migrate-up:
	@docker-compose exec app migrate -path /app/migrations -database "postgres://postgres:postgres@postgres:5432/restaurant_db?sslmode=disable" up

docker-migrate-down:
	@docker-compose exec app migrate -path /app/migrations -database "postgres://postgres:postgres@postgres:5432/restaurant_db?sslmode=disable" down 1

swagger:
	@swag init -g $(MAIN_PATH) -o ./docs

lint:
	@golangci-lint run ./...