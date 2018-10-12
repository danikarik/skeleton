HAS_REFLEX := $(shell command -v reflex;)

all: prep key cert build

help: ## Show usage
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

vendor: ## Install development tools
ifndef HAS_REFLEX
	@sh -c 'GO111MODULE=off go get -u -v github.com/cespare/reflex'
else
	@echo 'Already installed.'
endif

prep: ## Prepare project structure
	@mkdir -p certs

key: ## Generate private key
	@echo 'Generating private key...'
	@openssl genrsa -out certs/server.key 4096

cert: key ## Generate SSL certificate
	@echo 'Generating SSL certificate...'
	@openssl req -new -x509 -key certs/server.key -out certs/server.cert -days 3650 -subj /CN=localhost

build: ## Build binaries
	@echo 'Building server binary...'
	@go build -o bin/server ./cmd/server

serve: ## Run server
	@echo 'Starting server...'
	@bin/server

watch: ## Hot reload
	@reflex -r '\.go$$' -s -- sh -c 'make build && make serve'
