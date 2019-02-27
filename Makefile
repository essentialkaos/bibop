################################################################################

.DEFAULT_GOAL := help
.PHONY = fmt all clean deps gen-fuzz deps-test test help

################################################################################

all: bibop ## Build all binaries

bibop: ## Build bibop binary
	go build bibop.go

install: ## Install binaries
	cp bibop /usr/bin/bibop

uninstall: ## Uninstall binaries
	rm -f /usr/bin/bibop

deps: ## Download dependencies
	git config --global http.https://pkg.re.followRedirects true
	go get -d -v pkg.re/essentialkaos/ek.v10

deps-test: ## Download dependencies for tests
	git config --global http.https://pkg.re.followRedirects true
	go get -d -v pkg.re/check.v1

test: ## Run tests
	go test -covermode=count ./parser ./recipe

fmt: ## Format source code with gofmt
	find . -name "*.go" -exec gofmt -s -w {} \;

gen-fuzz: ## Generate archives for fuzz testing
	go-fuzz-build -o parser-fuzz.zip github.com/essentialkaos/bibop/parser

clean: ## Remove generated files
	rm -f bibop

help: ## Show this info
	@echo -e '\nSupported targets:\n'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[33m%-12s\033[0m %s\n", $$1, $$2}'
	@echo -e ''

################################################################################
