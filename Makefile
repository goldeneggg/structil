PKG_STRUCTIL := github.com/goldeneggg/structil
PKG_DYNAMICSTRUCT := github.com/goldeneggg/structil/dynamicstruct

PROFDIR := ./.prof
BENCH_RESULT_OLD := $(PROFDIR)/bench.old
BENCH_RESULT_NEW := $(PROFDIR)/bench.new
TRACE := $(PROFDIR)/trace.out
TESTBIN_STRUCTIL := $(PROFDIR)/structil.test
TESTBIN_DYNAMICSTRUCT := $(PROFDIR)/dynamicstruct.test

SRCS = $(shell find . -type f -name '*.go' | \grep -v 'vendor')
PKGS = $(shell ./scripts/_packages.sh)
TOOL_PKGS = $(shell ./scripts/_tool_packages.sh)

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

.PHONY: test
test:
	@go test -race -cover -parallel 2 $(PKGS)

.PHONY: -mk-profdir
-mk-profdir:
	@[ -d $(PROFDIR) ] || mkdir $(PROFDIR)

.PHONY: -mv-bench-result
-mv-bench-result:
	@[ ! -f $(BENCH_RESULT_NEW) ] || mv $(BENCH_RESULT_NEW) $(BENCH_RESULT_OLD)

benchmark = GOMAXPROCS=1 go test -run=NONE -bench . -benchmem -benchtime=100ms $1 $2 | tee $(BENCH_RESULT_NEW)

.PHONY: bench
bench: -mk-profdir -mv-bench-result
	@$(call benchmark,,$(PKGS))

.PHONY: show-latest-bench
show-latest-bench:
	@cat $(BENCH_RESULT_NEW)

.PHONY: benchstat
benchstat:
	@benchstat $(BENCH_RESULT_OLD) $(BENCH_RESULT_NEW)

benchmark-pprof = $(call benchmark,-cpuprofile $(PROFDIR)/$1.cpu.out -memprofile $(PROFDIR)/$1.mem.out -o $(PROFDIR)/$1.test,$2)

.PHONY: bench-prof
bench-prof: -mk-profdir -mv-bench-result
	@for pkg in $(PKGS); do echo ">>>>> Start: bench-prof for $${pkg}" && $(call benchmark-pprof,`basename $${pkg}`,$${pkg}); done

# pprof-cpu-structil OR pprof-cpu-dynamicstruct
.PHONY: pprof-cpu-%
pprof-cpu-%:
	@go tool pprof $(PROFDIR)/$*.test $(PROFDIR)/$*.cpu.out

# pprof-mem-structil OR pprof-mem-dynamicstruct
.PHONY: pprof-mem-%
pprof-mem-%:
	@go tool pprof $(PROFDIR)/$*.test $(PROFDIR)/$*.mem.out

test-trace: -mk-profdir
	@for pkg in $(PKGS); do echo ">>>>> Start: test-trace for $${pkg}" && go test -trace=$(PROFDIR)/`basename $${pkg}`.trace.out -o $(PROFDIR)/`basename $${pkg}`.test $${pkg}; done

.PHONY: trace-%
trace-%:
	@go tool trace $(PROFDIR)/$*.test $(PROFDIR)/$*.trace.out

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

lint-travis:
	@travis lint --org --debug .travis.yml

.PHONY: godoc
godoc:
	@godoc -http=:6060

mod-dl:
	@GO111MODULE=on go mod download

mod-tidy:
	@GO111MODULE=on go mod tidy

# Note: tools additional process as follows
#  - Add pacakge into tools.go
#  - Run "make mod-tidy"
#  - Run "make mod-tool-install"
mod-tools-install:
	@GO111MODULE=on go install $(TOOL_PKGS)

mod-golint-install:
	@GO111MODULE=on go install golang.org/x/lint/golint

mod-benchstat-install:
	@GO111MODULE=on go install golang.org/x/perf/cmd/benchstat

.PHONY: vendor
vendor:
	@GO111MODULE=on go mod vendor

.PHONY: clean
clean:
	@go clean -i -x -cache -testcache $(PKGS) $(TOOL_PKGS)
	rm -f $(BENCH_RESULT_OLD)
	rm -f $(BENCH_RESULT_NEW)
	rm -f $(PROFDIR)/*.test
	rm -f $(PROFDIR)/*.out

# CAUTION: this target removes all mod-caches
.PHONY: clean-mod-cache
clean-mod-cache:
	@go clean -i -x -modcache $(PKGS) $(TOOL_PKGS)
