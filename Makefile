GO ?= go
TEST_PACKAGE = "./..."

test: lint gotestsum goverreport
	@echo "running unit test..."
	@mkdir -p output
	$(GOTESTSUM) --format=pkgname --jsonfile=./output/out.json --packages=$(TEST_PACKAGE) -- -race -covermode=atomic -coverprofile=output/coverage.out -coverpkg $(TEST_PACKAGE)
	$(GOVERREPORT) -coverprofile=./output/coverage.out

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
