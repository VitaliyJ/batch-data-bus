docs:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g ./internal/server/http/router.go -o .
	rm docs.go swagger.json
	mv swagger.yaml openapi.yaml
	go mod tidy

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.0
	go mod tidy
	golangci-lint run --config .golangci.yml ./...
