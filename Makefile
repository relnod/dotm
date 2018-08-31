# This Makefile defines a set of rules to help in the development proccess.
# Run `make` to run all tests and checks on the project.
# For development you can run `make dev`.
#
# You can run `make help RULE={rule}` to get specific info about a rule.
# Example:
#    make help RULE=verify
####

export ROOT ?= $(shell pwd)
export TMP_ROOT ?= /tmp/dotm_tmp
export BASE_PACKAGE ?= github.com/relnod/dotm

GOARCH ?= amd64
TARGET ?= linux

# === all ===
# Runs all tests and verify scripts
.PHONY: all
all: test

# === test ===
# Runs all unit and e2e tests
.PHONY: test
test: test-unit test-e2e

# === test ===
# Runs all unit tests
.PHONY: test
test-unit:
	@echo "Running unit tests"
	@echo ""
	go test -v `go list ./... | grep -v test/`

# === test ===
# Runs all e2e tests
.PHONY: test-e2e
test-e2e:
	@echo "Running e2e tests"
	@echo ""
	go test -v ./test/e2e/...

# === update ===
# Updates all generated files
.PHONY: update
update:
	./hack/make-rules/update.sh

# === verify ===
# Runs various verify checks
.PHONY: verify
verify:
	./hack/make-rules/verify.sh

# === dev ===
# Starts modd (runs tests on file change)
.PHONY: dev
dev:
	modd -f hack/dev/modd.conf

export RULE ?= help

build:
	CGO_ENABLED=0 GOOS="${TARGET}" GOARCH="${GOARCH}" go build

# === help ===
# Prints this help message
.PHONY: help
help:
	@./hack/make-rules/help.sh
