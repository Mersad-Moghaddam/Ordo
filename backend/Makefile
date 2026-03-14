GO_PACKAGES=$(shell go list ./...)

.PHONY: run test revive migrate sqlc tidy benchmark

run:
	go run ./cmd/api

test:
	go test ./... -coverprofile=coverage.out

revive:
	revive -config revive.toml ./...

migrate:
	migrate -path migrations -database "$${ORDO_MYSQL_MIGRATE_DSN}" up

sqlc:
	sqlc generate

tidy:
	go mod tidy

benchmark:
	go test ./internal/usecase/platform -bench=. -benchmem
