run:
	go run main.go

mg-up:
	migrate -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -path migrations up

mg-down:
	migrate -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -path migrations down

mg-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

test:
	go test ./...

dc-up:
	docker compose up -d

.PHONY: run mg-up mg-down mg-install test dc-up