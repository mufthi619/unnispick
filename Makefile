include .env

.PHONY: migrate-create migrate-up migrate-down swagger

POSTGRES_URL=postgres://hepytech:pass2login@localhost:5432/latest_ok?sslmode=disable

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

migrate-up:
	migrate -path migrations -database "${POSTGRES_URL}" up

migrate-down:
	migrate -path migrations -database "${POSTGRES_URL}" down

swagger:
	swag init -g docs/main.go -o docs