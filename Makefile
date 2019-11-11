
EXAMPLE_SPEC_PATH ?= example.yml
EXAMPLE_APP_PATH ?= github.com/68696c6c/capricorn-example

.PHONY: image dep cli local-down test migrate

.DEFAULT:
	@echo 'Invalid target.'
	@echo
	@echo '    deps                                install dependancies'
	@echo '    build                               build the CLI for the current machine'
	@echo '    test                                run unit tests'
	@echo '    SPEC_PATH=/full/path/to/spec new    generate a new Goat project; provide SPEC_PATH'
	@echo '    example                             generate an example Goat project'
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
	capricorn new $(SPEC_PATH)

example: build
	rm -rf $(GOPATH)/src/$(EXAMPLE_APP_PATH)
	capricorn new $(EXAMPLE_SPEC_PATH)
	cd $(GOPATH)/src/$(EXAMPLE_APP_PATH) && make local
