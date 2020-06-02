PKG_STRUCTIL := github.com/goldeneggg/structil
PKG_DYNAMICSTRUCT := github.com/goldeneggg/structil/dynamicstruct

TESTDIR := ./.test
BENCH_OLD := $(TESTDIR)/bench.old
BENCH_NEW := $(TESTDIR)/bench.new
TRACE := $(TESTDIR)/trace.out
TESTBIN_STRUCTIL := $(TESTDIR)/structil.test
TESTBIN_DYNAMICSTRUCT := $(TESTDIR)/dynamicstruct.test

SRCS = $(shell find . -type f -name '*.go' | \grep -v 'vendor')
PKGS = $(shell ./scripts/_packages.sh)
TOOL_PKGS = $(shell cat ./tools/tools.go | grep _ | awk -F'"' '{print $$2}')

.DEFAULT_GOAL := test

.PHONY: version
version:
	@echo $(shell ./scripts/_version.sh)

.PHONY: pkgs
pkgs:
	@echo $(PKGS)

.PHONY: tool-pkgs
tool-pkgs:
	@echo $(TOOL_PKGS)

mod-dl:
	@GO111MODULE=on go mod download

mod-tidy:
	@GO111MODULE=on go mod tidy

# Note: tools additional process as follows
#  - Add pacakge into tools.go
#  - Run "make mod-tidy"
#  - Run "make mod-tools-install"
mod-tools-install: mod-tidy
	@GO111MODULE=on go install $(TOOL_PKGS)

mod-golint-install: mod-tidy
	@GO111MODULE=on go install golang.org/x/lint/golint

mod-benchstat-install: mod-tidy
	@GO111MODULE=on go install golang.org/x/perf/cmd/benchstat

.PHONY: test
test:
	@go test -race -cover -parallel 2 $(PKGS)

.PHONY: lint
lint: mod-golint-install
	@golint -set_exit_status $(PKGS)

.PHONY: vet
vet:
	@go vet $(PKGS)

ci-test:
	@./scripts/ci-test.sh

.PHONY: ci
ci: ci-test vet lint

.PHONY: -mk-testdir
-mk-testdir:
	@[ -d $(TESTDIR) ] || mkdir $(TESTDIR)

.PHONY: -mv-bench-result
-mv-bench-result:
	@[ ! -f $(BENCH_NEW) ] || mv $(BENCH_NEW) $(BENCH_OLD)

benchmark = go test -run=NONE -bench . -benchmem -cpu 1,2 -benchtime=500ms -count=5 $1 $2 | tee $(BENCH_NEW)

.PHONY: bench
bench: -mk-testdir -mv-bench-result
	@$(call benchmark,,$(PKGS))

.PHONY: benchstat
benchstat: mod-benchstat-install $(BENCH_OLD) $(BENCH_NEW)
	@benchstat $(BENCH_OLD) $(BENCH_NEW)

# WIP
.PHONY: benchstat-gist
benchstat-gist: mod-benchstat-install
	@bash -c "benchstat <(curl -sSL $${BENCH_OLD_GIST_URL}) <(curl -sSL $${BENCH_NEW_GIST_URL})"

benchmark-pprof = $(call benchmark,-cpuprofile $(TESTDIR)/$1.cpu.out -memprofile $(TESTDIR)/$1.mem.out -o $(TESTDIR)/$1.test,$2)

.PHONY: bench-prof
bench-prof: -mk-testdir -mv-bench-result
	@for pkg in $(PKGS); do echo ">>>>> Start: bench-prof for $${pkg}" && $(call benchmark-pprof,`basename $${pkg}`,$${pkg}); done

# pprof-cpu-structil OR pprof-cpu-dynamicstruct
.PHONY: pprof-cpu-%
pprof-cpu-%:
	@go tool pprof $(TESTDIR)/$*.test $(TESTDIR)/$*.cpu.out

# pprof-mem-structil OR pprof-mem-dynamicstruct
.PHONY: pprof-mem-%
pprof-mem-%:
	@go tool pprof $(TESTDIR)/$*.test $(TESTDIR)/$*.mem.out

test-trace: -mk-testdir
	@for pkg in $(PKGS); do echo ">>>>> Start: test-trace for $${pkg}" && go test -trace=$(TESTDIR)/`basename $${pkg}`.trace.out -o $(TESTDIR)/`basename $${pkg}`.test $${pkg}; done

.PHONY: trace-%
trace-%:
	@go tool trace $(TESTDIR)/$*.test $(TESTDIR)/$*.trace.out

.PHONY: godoc
godoc:
	@godoc -http=:6060

.PHONY: vendor
vendor:
	@GO111MODULE=on go mod vendor

.PHONY: clean
clean:
	@go clean -i -x -cache -testcache $(PKGS) $(TOOL_PKGS)
	rm -f $(BENCH_OLD)
	rm -f $(BENCH_NEW)
	rm -f $(TESTDIR)/*.test
	rm -f $(TESTDIR)/*.out

# CAUTION: this target removes all mod-caches
.PHONY: clean-mod-cache
clean-mod-cache:
	@go clean -i -x -modcache $(PKGS) $(TOOL_PKGS)


#####
#
# for Docker
#
#####

DOCKER_DIR := ./docker
DOCKER_IMAGE_MOD := structil/mod
DOCKER_IMAGE_TEST := structil/test

# PENDING
-docker-build-for-mod:
	@docker image build -t $(DOCKER_IMAGE_MOD) -f $(DOCKER_DIR)/mod/Dockerfile .

# -docker-build-for-test: -docker-build-for-mod
-docker-build-for-test:
	@docker image build -t $(DOCKER_IMAGE_TEST) -f $(DOCKER_DIR)/test/Dockerfile .

docker-test: -docker-build-for-test
	@docker container run --rm --cpus 2 $(DOCKER_IMAGE_TEST) test

docker-lint: -docker-build-for-test
	@docker container run --rm $(DOCKER_IMAGE_TEST) lint

docker-bench: -docker-build-for-test
	@docker container run --rm --cpus 2 -v `pwd`/.test:/go/src/github.com/goldeneggg/structil/.test:cached $(DOCKER_IMAGE_TEST) bench

hadolint: 
	@hadolint Dockerfile
