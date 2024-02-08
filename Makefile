include .env

mg-create:
	@migrate create -ext sql -dir db/migration -seq $(name)

mg-up:
	@migrate -path db/migration -database $(POSTGRES_URL) --verbose up

mg-down:
	@migrate -path db/migration -database $(POSTGRES_URL) --verbose down

seed:
	@go run cmd/seed/main.go

migrate:
	@go run cmd/migrate/main.go

build:
	@sqlc generate
	@go build

dev:
	@sqlc generate
	@air .
