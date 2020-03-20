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

benchmark-pprof = $(call benchmark,-cpuprofile $(PROFDIR)/$1.cpu.out -memprofile $(PROFDIR)/$1.mem.out -o $(PROFDIR)/$1.test,$2)

.PHONY: bench-prof
bench-prof: -mk-profdir -mv-bench-result
	@for pkg in $(PKGS); do echo ">>>>> Start: bench-prof for $${pkg}" && $(call benchmark-pprof,`basename $${pkg}`,$${pkg}); done

pprof = go tool pprof $1 $2

# pprof-cpu-structil OR pprof-cpu-dynamicstruct
.PHONY: pprof-cpu-%
pprof-cpu-%:
	@$(call pprof,$(PROFDIR)/$*.test,$(PROFDIR)/$*.cpu.out)

# pprof-mem-structil OR pprof-mem-dynamicstruct
.PHONY: pprof-mem-%
pprof-mem-%:
	@$(call pprof,$(PROFDIR)/$*.test,$(PROFDIR)/$*.mem.out)

test-trace = go test -trace=$(PROFDIR)/$1.trace.out -o $(PROFDIR)/$1.test $2

test-trace: -mk-profdir
	@for pkg in $(PKGS); do echo ">>>>> Start: test-trace for $${pkg}" && $(call test-trace,`basename $${pkg}`,$${pkg}); done

trace = go tool trace $(PROFDIR)/$1.test $(PROFDIR)/$1.trace.out

.PHONY: trace-%
trace-%:
	@$(call trace,$*)

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
