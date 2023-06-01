include $(PWD)/.env

haha:
	@echo ${OSFLAG}

.PHONY: build
build:
	$(GO) build $(LDFLAGS)

.PHONY: release
release: export BUILD = release
release: build

.PHONY: test
test: lint gotestsum goverreport prepare-envtest
	@echo "running unit test..."
	@mkdir -p output
	. $(PWD)/bin/testbin/test.env; \
		$(GOTESTSUM) --format=pkgname --jsonfile=./output/out.json --packages=$(TEST_PACKAGE) -- -race -covermode=atomic -coverprofile=output/coverage.out -coverpkg $(TEST_PACKAGE) $(LDFLAGS)
	$(GOVERREPORT) -coverprofile=./output/coverage.out

.PHONY: prepare-envtest
prepare-envtest: setup-envtest
	@# Prepare a k8s testenv
	@mkdir -p $(PWD)/bin/testbin
	@$(SETUP_ENVTEST) use --bin-dir "$(PWD)/bin/testbin" 1.27.1 -v debug -p env > $(PWD)/bin/testbin/test.env

SETUP_ENVTEST = $(shell pwd)/bin/setup-envtest
setup-envtest:
	$(call go_get,$(SETUP_ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest,latest)

GOLANGCI_LINT_FLAGS=
ifneq ($(GOMAXPROCS),)
GOLANGCI_LINT_FLAGS+=--concurrency=$(GOMAXPROCS)
endif
.PHONY: lint
lint: golangci_lint
	@echo ">> linting code..."
	$(GOLANGCI_LINT) $(GOLANGCI_LINT_FLAGS) run

GOTESTSUM = $(shell pwd)/bin/gotestsum
gotestsum:
	$(call go_get,$(GOTESTSUM),gotest.tools/gotestsum,v1.10.0)

GOVERREPORT = $(shell pwd)/bin/goverreport
goverreport:
	$(call go_get,$(GOVERREPORT),github.com/mcubik/goverreport,v1.0.0)

GOLANGCI_LINT = $(shell pwd)/bin/golangci-lint
golangci_lint:
	$(call go_get,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,v1.52.2)

# go_get will 'go get' any package $2@$3 and install it to $1.
# usage $(call go_get,$BinaryLocalPath,$GoModuleName,$Version)
define go_get
@set -e; \
if [ -f ${1} ]; then \
	[ -z ${3} ] && exit 0; \
	installed_version=$$(go version -m "${1}" | grep -E '[[:space:]]+mod[[:space:]]+' | awk '{print $$3}') ; \
	[ "$${installed_version}" = "${3}" ] && exit 0; \
	echo ">> ${1} ${2} $${installed_version}!=${3}, ${3} will be installed."; \
fi; \
module=${2}; \
if ! [ -z ${3} ]; then module=${2}@${3}; fi; \
echo "Downloading $${module}" ;\
GOBIN=$(shell pwd)/bin $(GO) install $${module} ;
endef
