upfile: main.go
	@go build

.PHONY: fmt
fmt:
	@golangci-lint run --fix

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: test
test:
	@go test ./...
