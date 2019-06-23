SRCS = $(shell find . -type f -name '*.go' | \grep -v 'vendor')
PKGS = $(shell go list ./... | \grep -v 'vendor')

.PHONY: pkgs
pkgs:
	@echo $(PKGS)

.PHONY: test
test:
	@go test -race -cover -parallel 2 $(PKGS)

.PHONY: test-v
test-v:
	@go test -v -race -cover -parallel 2 $(PKGS)

.PHONY: bench
bench:
	@GOMAXPROCS=1 go test -bench . -benchmem -benchtime=1s $(PKGS)

.PHONY: bench-v
bench-v:
	@GOMAXPROCS=1 go test -v -bench . -benchmem -benchtime=1s $(PKGS)

.PHONY: lint
lint:
	@golint -set_exit_status $(PKGS)

.PHONY: vet
vet:
	@go vet $(PKGS)

ci-test:
	@./scripts/ci-test.sh

.PHONY: ci
ci: ci-test vet lint

.PHONY: golangci
golangci:
	@golangci-lint run -c .golangci.yml

lint-travis:
	@travis lint --org --debug .travis.yml

.PHONY: godoc
godoc:
	@godoc -http=:6060

mod-dl:
	@GO111MODULE=on go mod download

mod-tidy:
	@GO111MODULE=on go mod tidy

.PHONY: vendor
vendor:
	@GO111MODULE=on go mod vendor
