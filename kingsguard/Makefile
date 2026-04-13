.PHONY: generate run test lint db schema

DB_NAME=kingsguard
DB_USER=postgres

generate:
	go install tool
	buf generate

run:
	go run cmd/main.go --config config/local.yaml

test:
	go test ./...

lint:
	golangci-lint run

db:
	psql postgres -c "drop database if exists $(DB_NAME);" -U $(DB_USER)
	createdb $(DB_NAME) -E 'utf8' -l 'en_US.UTF-8' -T 'template0' -U $(DB_USER)
	goose -allow-missing -dir migrations postgres "user=$(DB_USER) dbname=$(DB_NAME) sslmode=disable" up
	make schema


schema:
	pg_dump -d $(DB_NAME) -U ${DB_USER} --schema-only --no-owner --no-privileges --no-tablespaces --no-security-labels > schema.sql