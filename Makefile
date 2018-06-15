.PHONY: test
test:
	go test -v ./...

.PHONY: generate
generate:
	./hack/generate.sh

.PHONY: dev
dev:
	modd
