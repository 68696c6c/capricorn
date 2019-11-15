package ops

import (
	"github.com/68696c6c/capricorn/generator/models"
	"github.com/68696c6c/capricorn/generator/utils"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const dockerfileTemplate = `
FROM golang:1.13-alpine as env

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor
ENV CGO_ENABLED=0
ENV GOPROXY=https://proxy.golang.org,direct
# unfortunate, but needed if it falls back to direct, can't exclude a particular package
ENV GOSUMDB=off

RUN apk add --no-cache git gcc python bash openssh mysql-client

RUN mkdir -p /go/src/bitbucket.org/clearlink/loom-build
WORKDIR /go/src/bitbucket.org/clearlink/loom-build

RUN wget https://github.com/go-swagger/go-swagger/releases/download/v0.19.0/swagger_linux_amd64 -O /usr/local/bin/swagger
RUN chmod +x /usr/local/bin/swagger

RUN git config --global url."git@bitbucket.org:".insteadOf https://bitbucket.org/


################################################################################
# Local development stage.
FROM env as dev
RUN GOFLAGS="" go get -u github.com/go-delve/delve/cmd/dlv
RUN echo 'alias ll="ls -lah"' >> ~/.bashrc
`

const dockerComposeTemplate = `
version: "3.5"

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
      - ./:/go/src/capricorn-test
      - $HOME/.ssh:/root/.ssh:ro
    working_dir: /go/src/capricorn-test
    ports:
      - "80"
    env_file:
      - .app.env
    #  - .secret.env
    environment:
      VIRTUAL_HOST: capricorn.local
      ENV: local
      HTTP_PORT: 80
      MIGRATION_PATH: /go/src/capricorn-test/src/database
    networks:
      default:
        aliases:
          - capricorn.local

  db:
    image: mysql:5.7
    environment:
      MYSQL_ROOT_PASSWORD: secret
      MYSQL_DATABASE: build
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
DB_DATABASE=build
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

	return nil
}
