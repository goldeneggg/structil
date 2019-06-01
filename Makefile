SRCS = $(shell find . -type f -name '*.go' | \grep -v 'vendor')
PKGS = $(shell go list ./... | \grep -v 'vendor')

.PHONY: pkgs
pkgs:
	@echo $(PKGS)

.PHONY: run
run:
	@go run examples/misc/main.go

.PHONY: bench
bench:
	@go test -bench . $(PKGS)

mod-dl:
	@GO111MODULE=on go mod download

mod-tidy:
	@GO111MODULE=on go mod tidy

.PHONY: vendor
vendor:
	@GO111MODULE=on go mod vendor
