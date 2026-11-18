.DEFAULT_GOAL := help

.PHONY: help build test vet fmt fmt-check check doc

help:
	@printf '%s\n' \
		'make build      Compile all packages' \
		'make test       Run all tests' \
		'make vet        Run go vet' \
		'make fmt        Format Go files with gofmt' \
		'make fmt-check  Report files that need formatting' \
		'make check      Check formatting, run go vet, and run tests' \
		'make doc        Open local API documentation in a browser'

build:
	go build ./...

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

fmt-check:
	@files="$$(gofmt -l .)" || exit $$?; \
	if [ -n "$$files" ]; then \
		printf '%s\n' 'Run make fmt to format these files:' "$$files" >&2; \
		exit 1; \
	fi

check: fmt-check vet test

doc:
	go doc -http
