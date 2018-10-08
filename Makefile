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

TAG ?= latest

$(mkdir -p ./build)
$(mkdir -p ./artifacts)

# === all ===
# Runs all tests and verify scripts
.PHONY: all
all: verify test

# === test ===
# Runs all tests (unit and e2e)
.PHONY: test
test: test-unit test-e2e test-e2e-docker

# === test-unit ===
# Runs all unit tests
.PHONY: test-unit
test-unit:
	@echo "Running unit tests"
	@echo ""
	GO111MODULE=on go test -mod=vendor -v `go list ./... | grep -v /test`
	@echo ""

# === test-e2e ===
# Runs all e2e tests using the local binary
.PHONY: test-e2e
test-e2e: build
	@echo "Running e2e tests"
	@echo ""
	GO111MODULE=on go test -mod=vendor -v ${ROOT}/test -run=Normal
	@echo ""

# === test-e2e-docker ===
# Runs all e2e tests using the docker alias
.PHONY: test-e2e-docker
test-e2e-docker: build-docker
	@echo "Running e2e tests with docker alias"
	@echo ""
	GO111MODULE=on go test -mod=vendor -v ${ROOT}/test -run=Docker
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
	GO111MODULE=on CGO_ENABLED=0 GOOS="${TARGET}" GOARCH="${GOARCH}" go build -mod=vendor  -o ./build/dotm ./cmd/dotm

# === build-docker ===
# Builds the latest docker container
.PHONY: build-docker
build-docker:
	docker build . --tag=relnod/dotm:${TAG}

# === push-docker ===
# Builds and pushes the image to dockerhub
.PHONY: push-docker
push-docker: build-docker
	echo "docker login"
	@echo ${DOCKER_PWD} | docker login -u ${DOCKER_LOGIN} --password-stdin
	docker push relnod/dotm:${TAG}

# === release ===
# Creates a new release on github
.PHONY: release
release: release-artifacts
	echo "push release to github"
	@ghr -t ${GITHUB_TOKEN} -u relnod -r dotm -c ${COMMIT} -delete ${VERSION} ./artifacts/

# === release-artifacts ===
# Creates the release artifacts (in ./artifacts)
.PHONY: release-artifacts
release-artifacts:
	./scripts/make-rules/release-artifacts.sh

export RULE ?= help

# === help ===
# Prints this help message
.PHONY: help
help:
	@./scripts/make-rules/help.sh
