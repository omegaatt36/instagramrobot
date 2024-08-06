BUILD_DIR=build

.PHONY: fmt check test

all: dependency fmt check test

dependency:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/daixiang0/gci@latest
	go install github.com/client9/misspell/cmd/misspell@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

fmt:
	gofmt -s -w -l .
	@echo 'goimports' && goimports -w -local github.com/omegaatt36/instagramrobot $(shell find . -type f -name '*.go')
	gci write -s standard -s default -s "Prefix(github.com/omegaatt36/instagramrobot)" --skip-generated .
	go mod tidy
check:
	go vet -all ./...
	golangci-lint run
	misspell -error */**
	@echo 'staticcheck' && staticcheck $(shell go list ./...)

test:
	go test ./...