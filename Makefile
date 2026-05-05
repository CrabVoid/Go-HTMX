.PHONY: generate run test db-up db-down db-init

generate:
	templ generate
	sqlc generate

run:
	go run cmd/server/main.go

test:
	go test ./... -v

db-up:
	docker-compose up -d

db-down:
	docker-compose down

db-init:
	psql "postgres://postgres:postgres@localhost:5432/internship_manager?sslmode=disable" -f sql/schema.sql
