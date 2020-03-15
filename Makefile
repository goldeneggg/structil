PKG_STRUCTIL := github.com/goldeneggg/structil
PKG_DYNAMICSTRUCT := github.com/goldeneggg/structil/dynamicstruct

PROFDIR := ./.prof
BENCH_RESULT_OLD := $(PROFDIR)/bench.old
BENCH_RESULT_NEW := $(PROFDIR)/bench.new
CPU_PROF := $(PROFDIR)/cpu.out
MEM_PROF := $(PROFDIR)/mem.out
TRACE := $(PROFDIR)/trace.out
TESTBIN_STRUCTIL := $(PROFDIR)/structil.test
TESTBIN_DYNAMICSTRUCT := $(PROFDIR)/dynamicstruct.test

SRCS = $(shell find . -type f -name '*.go' | \grep -v 'vendor')
PKGS = $(shell go list ./... | \grep -v 'vendor')


.PHONY: pkgs
pkgs:
	@echo $(PKGS)

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

.PHONY: benchcmp
benchcmp:
	@benchcmp $(BENCH_RESULT_OLD) $(BENCH_RESULT_NEW)

benchmark-pprof = $(call benchmark,-cpuprofile $(CPU_PROF) -memprofile $(MEM_PROF),$1)

.PHONY: bench-prof-structil
bench-prof-structil: -mk-profdir -mv-bench-result
	@$(call benchmark-pprof,$(PKG_STRUCTIL))

.PHONY: bench-prof-dynamicstruct
bench-prof-dynamicstruct: -mk-profdir -mv-bench-result
	@$(call benchmark-pprof,$(PKG_DYNAMICSTRUCT))

pprof = go tool pprof $1

.PHONY: pprof-cpu
pprof-cpu:
	@$(call pprof,$(CPU_PROF))

.PHONY: pprof-mem
pprof-mem:
	@$(call pprof,$(MEM_PROF))

test-trace = go test -trace=$(TRACE) -o $1 $2

test-trace-structil: -mk-profdir
	@$(call test-trace,$(TESTBIN_STRUCTIL),$(PKG_STRUCTIL))

test-trace-dynamicstruct: -mk-profdir
	@$(call test-trace,$(TESTBIN_DYNAMICSTRUCT),$(PKG_DYNAMICSTRUCT))

trace = go tool trace $1 $(TRACE)

.PHONY: trace-structil
trace-structil:
	@$(call trace,$(TESTBIN_STRUCTIL))

.PHONY: trace-dynamicstruct
trace-dynamicstruct:
	@$(call trace,$(TESTBIN_DYNAMICSTRUCT))

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

.PHONY: vendor
vendor:
	@GO111MODULE=on go mod vendor
