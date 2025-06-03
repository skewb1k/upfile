.DEFAULT_GOAL := test

.PHONY: fmt
fmt:
	@go tool gofumpt -w -l .
	@go tool goimports -w -l .
	@go tool golangci-lint run --fix

.PHONY: lint
lint:
	@go tool golangci-lint run

.PHONY: test
test:
	@go tool gotestsum -f testdox

.PHONY: build
build:
	go build -o ~/.cache/bin/upfile ./cmd/upfile

.PHONY: gen
gen:
	@go generate ./...

