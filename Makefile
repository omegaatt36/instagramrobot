BUILD_DIR=build

.PHONY: fmt check test

all: fmt check test

fmt:
	gofmt -s -w -l .
	@echo 'goimports' && goimports -w -local gobe $(shell find . -type f -name '*.go' -not -path "./internal/*")
	gci write -s standard -s default -s "Prefix(gobe)" --skip-generated .
	go mod tidy

check:
	find . -name "*.json" | xargs -n 1 -t gojq . >/dev/null
	go vet -all ./...
	golangci-lint run
	misspell -error */**
	@echo 'staticcheck' && staticcheck $(shell go list ./... | grep -v internal)

vet:
	./vet.sh

test:
	go test ./...
