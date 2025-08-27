OS ?= $(shell uname | tr '[:upper:]' '[:lower:]')
ARCH ?= $(shell uname -m | tr '[:upper:]' '[:lower:]')
DATELOG := "[$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')]"
BINARY := build-scheduler


ifeq ($(ARCH),x86_64)
	ARCH=amd64
endif

.PHONY: default
default: help

.PHONY: build
build: ## Build the binary
	@mkdir -p $(CURDIR)/bin/$(OS)-$(ARCH)
	@echo "$(DATELOG) Building binary"
	GOOS=$(OS) GOARCH=$(ARCH) go build -o $(CURDIR)/bin/$(OS)-$(ARCH)/$(BINARY)
	@chmod +x $(CURDIR)/bin/$(OS)-$(ARCH)/$(BINARY)

.PHONY: run
run: ## Run the binary
	$(CURDIR)/bin/$(OS)-$(ARCH)/$(BINARY)

.PHONY: clean
clean: ## Clean /bin directory
	@rm -rf $(CURDIR)/bin

.PHONY: install
install: ## Install the binary using go install
	@echo "$(DATELOG) Installing $(BINARY)"
	GOOS=$(OS) GOARCH=$(ARCH) go install

.PHONY: lint
lint: ## Run golangci-lint
	@echo "$(DATELOG) Linting"
	golangci-lint run -v -c $(CURDIR)/.golangci.yml

.PHONY: test
test: ## Run go tests
	@echo "$(DATELOG) Running tests"
	go test ./...

.PHONY: tidy
tidy: ## Run go mod tidy
	@echo "$(DATELOG) Running go mod tidy"
	go mod tidy

.PHONY: vet
vet: ## Run go vet
	@echo "$(DATELOG) Running go vet"
	go vet ./...

.PHONY: help
help: ## Show this help
	@echo "Specify a command. The choices are:"
	@grep -hE '^[0-9a-zA-Z_-]+:.*?## .*$$' ${MAKEFILE_LIST} | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[0;36m%-20s\033[m %s\n", $$1, $$2}'
	@echo ""

