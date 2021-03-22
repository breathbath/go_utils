# set default shell
SHELL = bash -e -o pipefail

default: help

help:
	@echo "Usage: make [<target>]"
	@echo "where available targets are:"
	@echo
	@echo "help              : Print this help"
	@echo "test              : Run unit tests, if any"
	@echo "sca               : Run SCA"
	@echo "fmt               : Run gofmt and goimports"
	@echo

test:
	go test -race -v -p 1 ./...

sca:
	golangci-lint run

fmt:
	 goimports -w .
	 gofmt -w .
