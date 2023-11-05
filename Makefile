.PHONY: fmt
fmt: ## Run format tools on all go files
	gci write --skip-vendor --skip-generated \
        -s standard -s default -s "prefix(github.com/maypok86/benchmarks)" .
	gofumpt -l -w .

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL:= help
