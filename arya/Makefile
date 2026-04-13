.PHONY: proto_gen generate run test lint container mocks

proto_gen:
	go install tool
	buf generate --path vendor.protogen

gql:
	go run github.com/99designs/gqlgen generate

generate: proto_gen container mocks gql

run:
	go run cmd/main.go --config config/local.yaml

test:
	go test ./...

lint:
	golangci-lint run

container:
	go run scripts/clients_gen/main.go container > internal/config/container.go
	gofmt -w internal/config/container.go

mocks:
	go run scripts/clients_gen/main.go mocks