HAS_REFLEX := $(shell command -v reflex;)

all: help

help: ## Show usage
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

vendor: ## Install development tools
ifndef HAS_REFLEX
	@sh -c 'GO111MODULE=off go get -u -v github.com/cespare/reflex'
else
	@echo 'Already installed.'
endif

build: ## Build binaries
	@echo 'Building server binary...'
	@go build -o bin/server ./cmd/server

serve: ## Run server
	@echo 'Starting server...'
	@bin/server

watch: ## Hot reload
	@reflex -r '\.go$$' -s -- sh -c 'make build && make serve'
