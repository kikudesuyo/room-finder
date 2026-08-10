.PHONY: build dev lint lint-fix

GO111MODULE=on
LOAD_ENV=. ./default_env.sh
APP_BINARY ?= app

build:
	go build -o ./$(APP_BINARY) .

dev:
	$(LOAD_ENV) ; \
	go run main.go

lint:
	@type golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.4
	golangci-lint run ./...

lint-fix:
	golangci-lint run ./... --fix
