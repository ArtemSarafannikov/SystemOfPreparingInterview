.PHONY: proto_gen generate run test lint db testdb schema container mocks

DB_NAME=tirion
DB_USER=postgres

proto_gen:
	go install tool
	buf generate --path vendor.protogen
	buf generate --path api

generate: proto_gen container mocks

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

testdb:
	psql postgres -c "drop database if exists $(DB_NAME)_test;" -U $(DB_USER)
	createdb $(DB_NAME)_test -E 'utf8' -l 'en_US.UTF-8' -T 'template0' -U $(DB_USER)
	goose -allow-missing -dir migrations postgres "user=$(DB_USER) dbname=$(DB_NAME)_test sslmode=disable" up

schema:
	pg_dump -d $(DB_NAME) -U ${DB_USER} --schema-only --no-owner --no-privileges --no-tablespaces --no-security-labels > schema.sql

container:
	go run scripts/clients_gen/main.go container > internal/config/container.go
	gofmt -w internal/config/container.go

mocks:
	go run scripts/clients_gen/main.go mocks