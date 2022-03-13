PROJECT_NAME := "cloudflare-dyndns"

PKG := "github.com/rwxd/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/...)
GO_FILES := $(shell find . -name '*.go' | grep -v _test.go)

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

setup: ## Setup required things
	
	@python3 -m pip install -r requirements-dev.txt
	
	@go get -v -d ./...
	
