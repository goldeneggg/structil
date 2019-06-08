SRCS = $(shell find . -type f -name '*.go' | \grep -v 'vendor')
PKGS = $(shell go list ./... | \grep -v 'vendor')

.PHONY: pkgs
pkgs:
	@echo $(PKGS)

.PHONY: run
run:
	@go run examples/misc/main.go

.PHONY: test
test:
	@go test -race -cover -parallel 2 $(PKGS)

.PHONY: test-v
test-v:
	@go test -v -race -cover -parallel 2 $(PKGS)

.PHONY: bench
bench:
	@go test -bench . $(PKGS)

.PHONY: bench-v
bench-v:
	@go test -v -bench . $(PKGS)

.PHONY: lint
lint:
	@golint -set_exit_status $(PKGS)

.PHONY: vet
vet:
	@go vet $(PKGS)

ci-test:
	@./scripts/ci-test.sh

.PHONY: ci
ci: ci-test vet

mod-dl:
	@GO111MODULE=on go mod download

mod-tidy:
	@GO111MODULE=on go mod tidy

.PHONY: vendor
vendor:
	@GO111MODULE=on go mod vendor
