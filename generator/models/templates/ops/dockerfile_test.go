package ops

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDockerfile_MustParse(t *testing.T) {
	input := Dockerfile{
		Data: Ops{
			Workdir:      "example_workdir",
			AppHTTPAlias: "example_alias",
			MainDatabase: Database{
				Host:     "example_host",
				Port:     "1234",
				Username: "example_user",
				Password: "example_password",
				Name:     "example_dbname",
				Debug:    "1",
			},
		},
	}
	expected := `FROM golang:1.14-alpine as env

ENV CGO_ENABLED=0
ENV GOPROXY=https://proxy.golang.org,direct

RUN apk add --no-cache git gcc openssh mysql-client

RUN go get golang.org/x/tools/cmd/goimports

# Install swagger for generating API docs.
RUN go get -v github.com/go-swagger/go-swagger/cmd/swagger

# Install golangci-lint for linting.
RUN wget https://github.com/golangci/golangci-lint/releases/download/v1.24.0/golangci-lint-1.24.0-linux-amd64.tar.gz \
    && tar xzf golangci-lint-1.24.0-linux-amd64.tar.gz \
    && mv golangci-lint-1.24.0-linux-amd64/golangci-lint /usr/local/bin/golangci-lint

# Install goose for running migrations.
RUN go get -u github.com/pressly/goose/cmd/goose

RUN mkdir -p /example_workdir
WORKDIR /example_workdir


# Local development stage.
FROM env as dev

RUN go get -v github.com/go-delve/delve/cmd/dlv

RUN apk add --no-cache bash
RUN echo 'alias ll="ls -lah"' >> ~/.bashrc


# Pipeline stage for running unit tests, linters, etc.
FROM env as built

COPY ./src /example_workdir
RUN go build -i -o app
`
	result := input.MustParse()
	assert.Equal(t, expected, result)
}
