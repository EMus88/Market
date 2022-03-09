run:
	go run ./cmd/main.go
test:
	go test ./internal/handler/tests

swag:
	swag init -g ./cmd/main.go