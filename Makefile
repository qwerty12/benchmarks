.PHONY: setup
setup: deps ## Setup development environment
	cp ./scripts/pre-push.sh .git/hooks/pre-push
	chmod +x .git/hooks/pre-push

.PHONY: deps
deps: ## Install all the build and lint dependencies
	bash scripts/deps.sh

.PHONY: lint
lint: ## Run all the linters
	golangci-lint run -v ./...

.PHONY: ci
ci: lint ## Run all the tests and code checks

.PHONY: fmt
fmt: ## Run format tools on all go files
	gci write --skip-vendor --skip-generated \
        -s standard -s default -s "prefix(github.com/maypok86/benchmarks)" .
	gofumpt -l -w .

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL:= help

results:
	go run ./throughput/cmd/main.go "./throughput/results/throughput.txt"