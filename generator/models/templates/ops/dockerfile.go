package ops

import (
	"github.com/68696c6c/capricorn/generator/models/templates"
	"github.com/68696c6c/capricorn/generator/utils"
)

// @TODO use a hosted base image
const DockerfileTemplate = `FROM golang:1.14-alpine as env

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

RUN mkdir -p /{{ .Workdir }}
WORKDIR /{{ .Workdir }}


# Local development stage.
FROM env as dev

RUN go get -v github.com/go-delve/delve/cmd/dlv

RUN apk add --no-cache bash
RUN echo 'alias ll="ls -lah"' >> ~/.bashrc


# Pipeline stage for running unit tests, linters, etc.
FROM env as built

COPY ./src /{{ .Workdir }}
RUN go build -i -o app
`

type Dockerfile struct {
	Name templates.FileData `yaml:"name"`
	Path templates.PathData `yaml:"path"`

	Data Ops `yaml:"data"`
}

// This is only used for testing.
func (m Dockerfile) MustParse() string {
	result, err := utils.ParseTemplateToString("tmp_template_dockerfile", DockerfileTemplate, m.Data)
	if err != nil {
		panic(err)
	}
	return result
}

func (m Dockerfile) Generate() error {
	err := utils.GenerateFile(m.Path.Base, m.Name.Full, DockerfileTemplate, m.Data)
	if err != nil {
		return err
	}
	return nil
}
