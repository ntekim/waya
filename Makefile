.PHONY: run build sqlc swagger tidy

# Run the application
run:
	go run cmd/api/main.go

# Generate Go code from SQL
sqlc:
	sqlc generate

# Generate Swagger docs (requires swaggo)
swagger:
	swag init -g cmd/api/main.go --output docs

# Tidy up modules
tidy:
	go mod tidy

# Run all setup steps
setup: tidy sqlc swagger