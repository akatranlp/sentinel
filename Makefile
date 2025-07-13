MAKEFLAGS += --no-print-directory

.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  generate                          Generate"
	@echo "  examples/all                      Build all examples"
	@echo "  examples/basic                    Build basic example"
	@echo "  run/examples/basic                Run basic example"
	@echo "  setup                             Install dependencies"
	@echo "  clean                             Clean"
	@echo "  run/dev                           Run development environment"
	@echo "  dev/frontend                      Run frontend"
	@echo "  tidy                              Fmt and Tidy"
	@echo "  lint                              Lint"
	@echo "  test                              Test"
	@echo "  test/verbose                      Test verbose"
	@echo "  debug                             Debug"


.PHONY: run/dev
run/dev:
	@echo "Starting development environment..."
	@go tool concur || true

.PHONY: dev/frontend
dev/frontend:
	@echo "Starting development environment..."
	@cd web && pnpm dev || true

.PHONY: build/tailwind
build/tailwind:
	@echo "Executing tailwindcss..."
	@./node_modules/.bin/tailwindcss -i ./openid/web/assets/tailwind.css -o ./openid/web/assets/dist/styles.css --minify

.PHONY: setup
setup:
	@echo "Setting up environment..."
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

.PHONY: examples/all
examples/all: examples/basic

.PHONY: examples/basic
examples/basic: generate
	@echo "Building Basic Example"
	@go build -o ./bin/basic ./examples/basic

.PHONY: run/examples/basic
run/examples/basic: examples/basic
	@echo "Running basic example"
	./bin/basic

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

