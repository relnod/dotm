# This Makefile defines a set of rules to help in the development proccess.
# Run `make` to run all tests and checks on the project.
# For development you can run `make dev`.
#
# You can run `make help RULE={rule}` to get specific info about a rule.
# Example:
#    make help RULE=verify
####

export ROOT ?= $(shell pwd)/
export TMP_ROOT ?= /tmp/dotm_tmp
export BASE_PACKAGE ?= github.com/relnod/dotm

GOARCH ?= amd64
TARGET ?= linux

# === all ===
# Runs all tests and verify scripts
.PHONY: all
all: verify test

# === test ===
# Runs all tests (unit and e2e)
.PHONY: test
test: test-unit test-e2e

# === test-unit ===
# Runs all unit tests
.PHONY: test-unit
test-unit:
	@echo "Running unit tests"
	@echo ""
	GO111MODULE=on go test -mod=vendor -v `go list ./... | grep -v /test`
	@echo ""

# === test-e2e ===
# Runs all e2e tests
.PHONY: test-e2e
test-e2e: install
	@echo "Running e2e tests"
	@echo ""
	rm -rf test/testdata/remote/.git
	cd test/testdata/remote && \
		git init && \
	    git config --local user.email "you@example.com" && \
	    git config --local user.name "Your Name" && \
		git add . && \
		git commit -m "initial commit"
	GO111MODULE=on go test -mod=vendor -v ./test
	rm -rf test/testdata/remote/.git
	@echo ""

# === update ===
# Updates all generated files
.PHONY: update
update:
	./scripts/make-rules/update.sh

# === verify ===
# Runs various verify checks
.PHONY: verify
verify:
	./scripts/make-rules/verify.sh

# === watch ===
# Starts modd (runs tests on file change)
.PHONY: watch
watch:
	modd -f scripts/dev/modd.conf

# === install ===
# Installs dotm with go the go command
.PHONY: install
install:
	go install ./cmd/dotm

# === build ===
# Builds the dotm binary
.PHONY: build
build:
	CGO_ENABLED=0 GOOS="${TARGET}" GOARCH="${GOARCH}" go build ./cmd/dotm

export RULE ?= help

# === help ===
# Prints this help message
.PHONY: help
help:
	@./scripts/make-rules/help.sh
