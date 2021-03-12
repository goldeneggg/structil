LOCAL_GO := go$${LOCAL_GOVERSION}

PKG_STRUCTIL := github.com/goldeneggg/structil
PKG_DYNAMICSTRUCT := github.com/goldeneggg/structil/dynamicstruct

PKG_MAPSTRUCTURE := github.com/mitchellh/mapstructure
PKG_VIPER := github.com/spf13/viper
PKG_GOCMP := github.com/google/go-cmp

TESTDIR := ./.test
BENCH_OLD := $(TESTDIR)/bench.old
BENCH_NEW := $(TESTDIR)/bench.new
BENCH_LATEST_URL := https://raw.githubusercontent.com/goldeneggg/structil/bench-latest/BENCHMARK_LATEST.txt
TRACE := $(TESTDIR)/trace.out

SRCS = $(shell find . -type f -name '*.go' | \grep -v 'vendor')
PKGS = $(shell ./scripts/packages.sh)
TOOL_PKGS = $(shell cat ./tools/tools.go | grep _ | awk -F'"' '{print $$2}')

assert-command = $(if $(shell which $1),,$(error '$1' command is missing))


.DEFAULT_GOAL := test


###
# show informations
###
.PHONY: version
version:
	@echo $(shell ./scripts/version.sh)

.PHONY: pkgs
pkgs:
	@echo $(PKGS)

.PHONY: tool-pkgs
tool-pkgs:
	@echo $(TOOL_PKGS)

###
# manage modules
###
go-mod = GO111MODULE=on $(LOCAL_GO) mod $1 $2
go-install = GO111MODULE=on $(LOCAL_GO) install $1
chk_latest = go list -u -m $1

.PHONY: mod-dl
mod-dl:
	@$(call go-mod,download,)

.PHONY: mod-tidy
mod-tidy:
	@$(call go-mod,tidy,)

.PHONY: vendor
vendor:
	@$(call go-mod,vendor,)

# Note: tools additional process as follows
#  1. Add pacakge into tools.go
#  2. Run "make mod-tidy"
#  3. Run "make mod-tools-install"
.PHONY: mod-tools-install
mod-tools-install: mod-tidy
	@$(call go-install,$(TOOL_PKGS))

.PHONY: chk-latest-mapstructure
chk-latest-mapstructure:
	@$(call chk_latest,$(PKG_MAPSTRUCTURE))

.PHONY: chk-latest-viper
chk-latest-viper:
	@$(call chk_latest,$(PKG_VIPER))

.PHONY: chk-latest-gocmp
chk-latest-gocmp:
	@$(call chk_latest,$(PKG_GOCMP))

.PHONY: update-all-modules
update-all-modules:
	@go get -u && make test

upgrade_module = echo module-query="$2"; GO111MODULE=on go get $1@$2
upgrade_to_latest = $(call upgrade_module,$1,latest)

.PHONY: upgrade-latest-mapstructure
upgrade-latest-mapstructure:
	@$(call upgrade_to_latest,$(PKG_MAPSTRUCTURE))

.PHONY: upgrade-latest-viper
upgrade-latest-viper:
	@$(call upgrade_to_latest,$(PKG_VIPER))

.PHONY: upgrade-latest-gocmp
upgrade-latest-gocmp:
	@$(call upgrade_to_latest,$(PKG_GOCMP))


###
# run tests
###
run-test = $(LOCAL_GO) test -v -race -cover $1 $2

.PHONY: test
test:
	@$(call run-test,-p 4,$(PKGS))

.PHONY: subtest
subtest:
	@$(call run-test,-run,$(ST))

.PHONY: lint
lint: mod-golint-install
	@golint -set_exit_status $(PKGS)

.PHONY: vet
vet:
	@$(LOCAL_GO) vet $(PKGS)

.PHONY: ci-test
ci-test:
	@./scripts/ci-test.sh

.PHONY: ci
ci: ci-test vet lint

###
# run benchmark and profile
###
.PHONY: -mk-testdir
-mk-testdir:
	@[ -d $(TESTDIR) ] || mkdir $(TESTDIR)

.PHONY: -mv-bench-result
-mv-bench-result:
	@[ ! -f $(BENCH_NEW) ] || mv $(BENCH_NEW) $(BENCH_OLD)

benchmark = $(LOCAL_GO) test -run=NONE -bench . -benchmem -cpu 1,2 -benchtime=500ms -count=5 $1 $2 | tee $(BENCH_NEW)

.PHONY: bench
bench: -mk-testdir -mv-bench-result
	@$(call benchmark,,$(PKGS))

.PHONY: benchstat
benchstat: mod-benchstat-install $(BENCH_OLD) $(BENCH_NEW)
	@benchstat $(BENCH_OLD) $(BENCH_NEW)

.PHONY: benchstat-ci
benchstat-ci: mod-benchstat-install
	@bash -c "benchstat <(curl -sSL $(BENCH_LATEST_URL)) $(BENCH_NEW)"

benchmark-pprof = $(call benchmark,-cpuprofile $(TESTDIR)/$1.cpu.out -memprofile $(TESTDIR)/$1.mem.out -o $(TESTDIR)/$1.test,$2)

.PHONY: bench-prof
bench-prof: -mk-testdir -mv-bench-result
	@for pkg in $(PKGS); do echo ">>>>> Start: bench-prof for $${pkg}" && $(call benchmark-pprof,`basename $${pkg}`,$${pkg}); done

# pprof-cpu-structil OR pprof-cpu-dynamicstruct
.PHONY: pprof-cpu-%
pprof-cpu-%:
	@$(LOCAL_GO) tool pprof $(TESTDIR)/$*.test $(TESTDIR)/$*.cpu.out

# pprof-mem-structil OR pprof-mem-dynamicstruct
.PHONY: pprof-mem-%
pprof-mem-%:
	@$(LOCAL_GO) tool pprof $(TESTDIR)/$*.test $(TESTDIR)/$*.mem.out

.PHONY: test-trace
test-trace: -mk-testdir
	@for pkg in $(PKGS); do echo ">>>>> Start: test-trace for $${pkg}" && $(LOCAL_GO) test -trace=$(TESTDIR)/`basename $${pkg}`.trace.out -o $(TESTDIR)/`basename $${pkg}`.test $${pkg}; done

.PHONY: trace-%
trace-%:
	@$(LOCAL_GO) tool trace $(TESTDIR)/$*.test $(TESTDIR)/$*.trace.out

###
# use local specified go version
# See: https://pkg.go.dev/golang.org/dl
#
# [Usage with direnv]
# 1. Write contents as follows into .envrc using "direnv edit" and save.
#
# >>>>>>>>>>
# # Setup specified local go version
# export LOCAL_GOVERSION=
#
# setup_specified_go_version() {
#   go get golang.org/dl/go${LOCAL_GOVERSION}
#   go${LOCAL_GOVERSION} download
#   echo "Use go ${LOCAL_GOVERSION}"
# }
#
# if [ ! -z ${LOCAL_GOVERSION} ]
# then
#   setup_specified_go_version
# fi
# <<<<<<<<<<
#
# 2. Confirm local go version with "make show-local-go" command. (should be printed "go" )
# 3. If you want to switch go version to specified number, then open .envrc and write "export LOCAL_GOVERSION=<YOUR_VERSION>" and save.
# 4. After save, setup process for specified go version will be executed automatically by direnv mechanism.
# 5. Confirm NEW local go version with "make show-local-go" command. (should be printed "go<YOUR_VERSION>" )
# 6. Happy developing!!
#
###
.PHONY: show-local-go
show-local-go:
	@echo $(LOCAL_GO)

###
# clean up
###
.PHONY: clean
clean:
	@$(LOCAL_GO) clean -i -x -cache -testcache $(PKGS) $(TOOL_PKGS)
	rm -f $(BENCH_OLD)
	rm -f $(BENCH_NEW)
	rm -f $(TESTDIR)/*.test
	rm -f $(TESTDIR)/*.out

# CAUTION: this target removes all mod-caches
.PHONY: clean-mod-cache
clean-mod-cache:
	@$(LOCAL_GO) clean -i -x -modcache $(PKGS) $(TOOL_PKGS)

###
# miscellaneous
###
.PHONY: godoc
godoc:
	@godoc -http=:6060

#####
#
# with docker
#
#####

DOCKER_DIR := ./docker
DOCKER_IMAGE_MOD := structil/mod
DOCKER_IMAGE_TEST := structil/test

# PENDING
-docker-build-for-mod:
	@docker image build -t $(DOCKER_IMAGE_MOD) -f $(DOCKER_DIR)/mod/Dockerfile .

# docker-build-for-test: -docker-build-for-mod
.PHONY: docker-build-for-test
docker-build-for-test:
	@docker image build -t $(DOCKER_IMAGE_TEST) -f $(DOCKER_DIR)/test/Dockerfile .

.PHONY: docker-test
docker-test: docker-build-for-test
	@docker container run --rm --cpus 2 $(DOCKER_IMAGE_TEST) test

.PHONY: docker-lint
docker-lint: docker-build-for-test
	@docker container run --rm $(DOCKER_IMAGE_TEST) lint

.PHONY: docker-bench
docker-bench: docker-build-for-test
	@docker container run --rm --cpus 2 -v `pwd`/.test:/go/src/github.com/goldeneggg/structil/.test:cached $(DOCKER_IMAGE_TEST) bench

.PHONY: hadolint
hadolint: 
	@hadolint docker/**/Dockerfile*
