export GO111MODULE=on

all: test lint

test:
	go test ./...

lint:
	bash -c "diff -u <(echo -n) <(gofmt -d ./)"
	golint -set_exit_status $(go list ./... | grep -v /vendor/)

install:
	cd cmd/dotm && go install
