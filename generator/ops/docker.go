package ops

import (
	"github.com/68696c6c/capricorn/generator/models"
	"github.com/68696c6c/capricorn/generator/utils"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// @TODO just use a hosted image
const dockerfileTemplate = `FROM golang:1.14-alpine as env

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

const dockerComposeTemplate = `
version: "3.6"

networks:
  default:
    external:
      name: docker-dev

volumes:
  pkg:
  db-volume:

services:

  app:
    image: app:dev
    command: ./app server 80
    depends_on:
      - db
    volumes:
      - pkg:/go/pkg
      - ./:/{{ .Workdir }}
    working_dir: /{{ .Workdir }}
    ports:
      - "80"
    env_file:
      - .app.env
    environment:
      VIRTUAL_HOST: {{ .Workdir }}.local
      ENV: local
      HTTP_PORT: 80
      MIGRATION_PATH: /{{ .Workdir }}/src/database
    networks:
      default:
        aliases:
          - {{ .AppHTTPAlias }}.local

  db:
    image: mysql:5.7
    environment:
      MYSQL_ROOT_PASSWORD: secret
      MYSQL_DATABASE: {{ .AppDBName }}
    ports:
      - "${HOST_DB_PORT:-3310}:3306"
    volumes:
      - db-volume:/var/lib/mysql
`

const appEnvTemplate = `
DB_HOST=db
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=secret
DB_DATABASE={{ .AppDBName }}
DB_DEBUG=1
`

func CreateDocker(spec models.Project, logger *logrus.Logger) error {
	logPrefix := "CreateDocker | "
	logger.Debug(logPrefix, "generating docker files")

	err := utils.CreateDir(spec.Paths.Docker)
	if err != nil {
		return errors.Wrapf(err, "failed to create project directory '%s'", spec.Paths.Docker)
	}

	err = utils.GenerateFile(spec.Paths.Docker, "Dockerfile", dockerfileTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create Dockerfile")
	}

	err = utils.GenerateFile(spec.Paths.Root, "docker-compose.yml", dockerComposeTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create docker-compose.yml")
	}

	err = utils.GenerateFile(spec.Paths.Root, ".app.env", appEnvTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create .app.env")
	}

	err = utils.GenerateFile(spec.Paths.Root, ".app.example.env", appEnvTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create .app.env")
	}

	return nil
}
