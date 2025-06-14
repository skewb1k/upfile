.DEFAULT_GOAL := test


BINARY = upfile
ifeq ($(OS),Windows_NT)
	BINARY := $(BINARY).exe
endif


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
	go build -o dist/$(BINARY) ./cmd/upfile

.PHONY: clean
clean:
	rm -f dist/$(BINARY)
