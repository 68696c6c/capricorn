package ops

import (
	"os"

	"github.com/68696c6c/capricorn/generator/utils"
	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
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
    #env_file:
    #  - .app.env
    #  - .secret.env
    environment:
      VIRTUAL_HOST: capricorn.local
      ENV: debug
      LISTEN_PORT: 80
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

const makefileTemplate = `
NETWORK_NAME ?= docker-dev
APP_NAME = app
DB_NAME = db
MODULE = {{.Module}}

DOC_PATH_BASE = docs/swagger.json
DOC_PATH_FINAL = docs/api-spec.json

.PHONY: docs build

.DEFAULT:
	@echo 'App targets:'
	@echo
	@echo '    image              build the Docker image for local development'
	@echo '    build              build api image and compile the app'
	@echo '    deps               install dependancies'
	@echo '    setup-network      create local docker network'
	@echo '    setup              set up local databases'
	@echo '    local              spin up local environment'
	@echo '    local-down         tear down local environment'
	@echo '    test               run unit tests'
	@echo '    migrate            migrate the local database'
	@echo '    migration          create a new migration'
	@echo '    docs               build the Swagger docs'
	@echo '    docs-server        build and serve the Swagger docs'
	@echo '    lint               run the linter'
	@echo


default: .DEFAULT

image:
	docker build . -f docker/Dockerfile --target dev -t $(APP_NAME):dev

build:
	docker-compose run --rm $(APP_NAME) go build -i -o $(APP_NAME)

init:
	docker-compose run --rm $(APP_NAME) go mod init $(MODULE)

deps:
	docker-compose run --rm $(APP_NAME) go mod tidy
	docker-compose run --rm $(APP_NAME) go mod vendor

setup-network:
	docker network create docker-dev

setup: image deps build
	docker-compose run --rm $(DB_NAME) mysql -u root -psecret -h $(DB_NAME) -e "CREATE DATABASE IF NOT EXISTS test_repos"
	docker-compose run --rm $(APP_NAME) bash -c "./$(APP_NAME) migrate install && ./$(APP_NAME) seed"

local: local-down build
	NETWORK_NAME="$(NETWORK_NAME)" docker-compose up

local-down:
	NETWORK_NAME="$(NETWORK_NAME)" docker-compose down

test:
	docker-compose run --rm $(APP_NAME) go test ./... -cover

migrate:
	docker-compose run --rm $(APP_NAME) ./$(APP_NAME) migrate

migration:
	docker-compose run --rm $(APP_NAME) ./$(APP_NAME) make:migration $(name)

docs: build
	docker-compose run --rm $(APP_NAME) bash -c "GO111MODULE=off swagger generate spec -mo '$(DOC_PATH_BASE)'"
	docker-compose run --rm $(APP_NAME) ./$(APP_NAME) gen:docs

docs-server: docs
	swagger serve "$(DOC_PATH_FINAL)"

lint:
	docker-compose run --rm $(APP_NAME) ./ops/scripts/lint.sh

`

func CreateDocker(spec utils.Spec, logger *logrus.Logger) error {
	logPrefix := "CreateDocker | "

	path := spec.Paths.Docker
	cwd, _ := os.Getwd()
	logger.Debug(logPrefix, "path: ", path, cwd)

	err := utils.CreateDir(path)
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

	err = utils.GenerateFile(spec.Paths.Root, "Makefile", makefileTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create Makefile")
	}

	return nil
}
