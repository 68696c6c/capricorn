
EXAMPLE_SPEC_PATH ?= example.yml
EXAMPLE_APP_PATH ?= ~/Code/Go/src/github.com/68696c6c/capricorn-test

.PHONY: image dep cli local-down test migrate

.DEFAULT:
	@echo 'Invalid target.'
	@echo
	@echo '    deps          install dependancies'
	@echo '    build         build the CLI for the current machine'
	@echo '    test          run unit tests'
	@echo '    new           generate a new Goat project'
	@echo

default: .DEFAULT

deps:
	go mod tidy
	go mod vendor

build:
	 go build -o /usr/local/bin/capricorn

test:
	go test ./... -cover

new: build
	capricorn new $(EXAMPLE_SPEC_PATH)
