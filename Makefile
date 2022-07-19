EXAMPLE_SPEC_PATH ?= loom-rnd.yml
EXAMPLE_APP_PATH ?= ~/Code


deps:
	go mod tidy
	go mod vendor

build:
	 go build -o capricorn

test:
	go test ./... -cover

example: build
	rm -rf ~/Code/example2
	./capricorn new $(EXAMPLE_SPEC_PATH) $(EXAMPLE_APP_PATH)

