export GO111MODULE=on

all: test lint

test:
	go test ./...

lint:
	bash -c "diff -u <(echo -n) <(gofmt -d ./)"
	golint -set_exit_status $(go list ./... | grep -v /vendor/)

install:
	cd cmd/dotm && go install

VERSION := $(shell dotm --version | cut -d ' ' -f 3)

# Make sure the GITHUB_TOKEN is set before running make release.
# After a release bump the version in cmd/dotm/commands/root.go.
release:
ifeq ($(GITHUB_TOKEN),)
	@echo "GITHUB_TOKEN not set"
	@exit
else
	git tag -a ${VERSION} -m "${VERSION}"
	git push --tags
	goreleaser --rm-dist
endif
