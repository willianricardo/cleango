# Go parameters
APP_NAME=cleango
MAIN=main.go

# Migration parameters
MIGRATE=migrate
DB_URL=sqlite3://$(shell pwd)/database/database.db
MIGRATIONS=database/migrations

# Swaggo
SWAG=swag

.PHONY: all build run test migrate-up migrate-down swag clean

all: build

build:
	go build -o $(APP_NAME) $(MAIN)

run:
	go run $(MAIN)

test:
	go test ./test/...

migrate-up:
	$(MIGRATE) -database $(DB_URL) -path $(MIGRATIONS) up

migrate-down:
	$(MIGRATE) -database $(DB_URL) -path $(MIGRATIONS) down

swag:
	$(SWAG) init --parseDependency --parseInternal

clean:
	rm -f $(APP_NAME)
	rm -f database/*.db