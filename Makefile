MAKEFLAGS += --no-print-directory
MAIN_PACKAGE_PATH := ./examples/basic/
BINARY_NAME := identity-provider
APP_VERSION ?= $(shell git describe --tags --always --dirty)
APP_GIT_COMMIT ?= $(shell git rev-parse HEAD)
APP_GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
APP_GIT_REPOSITORY ?= https://github.com/git-classrooms/identity-provider
APP_BUILD_TIME ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  generate                          Generate"
	@echo "  setup                             Install dependencies"
	@echo "  setup/ci                          Install dependencies for CI"
	@echo "  clean                             Clean"
	@echo "  run/dev                           Run development environment"
	@echo "  tidy                              Fmt and Tidy"
	@echo "  lint                              Lint"
	@echo "  test                              Test"
	@echo "  test/verbose                      Test verbose"
	@echo "  debug                             Debug"


.PHONY: run/dev
run/dev:
	@echo "Starting development environment..."
	@go tool concur || true

.PHONY: build/tailwind
build/tailwind:
	@echo "Executing tailwindcss..."
	@./node_modules/.bin/tailwindcss -i ./openid/web/assets/tailwind.css -o ./openid/web/assets/dist/styles.css --minify

.PHONY: setup
setup:
	@echo "Setting up environment..."
	@yarn install
	@go mod download

.PHONY: setup/ci
setup/ci:
	@echo "Installing..."
	@yarn install
	@go mod download

.PHONY: watch/tailwind
watch/tailwind:
	@echo "Executing tailwindcss..."
	@./node_modules/.bin/tailwindcss -i ./openid/web/assets/tailwind.css -o ./openid/web/assets/dist/styles.css --watch

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf ./bin ./tmp
	@rm -f  ./openid/web/assets/dist/styles.css
	@find . -type f -name '*_templ.go' -delete
	@find . -type f -name '*_enum.go' -delete

.PHONY: generate
generate:
	@echo "Generating code..."
	@go generate ./...
	@$(MAKE) build/tailwind

.PHONY: tidy
tidy:
	@echo "Tidying up..."
	go fmt ./...
	go tool templ fmt .
	go mod tidy

.PHONY: lint
lint:
	@echo "Linting..."
	go tool golangci-lint run

.PHONY: test
test:
	@echo "Testing..."
	go test -cover ./...

.PHONY: test/verbose
test/verbose:
	@echo "Testing..."
	go test -v -cover ./...

.PHONY: debug
debug:
	@echo "Debugging..."
	dlv debug

