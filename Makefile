.PHONY: build test lint run clean

# Variables
BINARY_NAME=news-app
MAIN_PATH=./cmd/server

# Build the application
build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

# Run tests
test:
	go test -v ./...


# Run the application
run:
	go run $(MAIN_PATH)

# Clean
clean:
	go clean
	rm -f $(BINARY_NAME)

# Run in Docker
docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down


# Run all checks
check: test
