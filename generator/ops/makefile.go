package ops

import (
	"github.com/68696c6c/capricorn/generator/models"
	"github.com/68696c6c/capricorn/generator/utils"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const makefileTemplate = `
NETWORK_NAME ?= docker-dev
APP_NAME = app
DB_NAME = db
MODULE = {{.Config.Module}}

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

func CreateMakefile(spec models.Project, logger *logrus.Logger) error {
	logPrefix := "CreateMakefile | "
	logger.Debug(logPrefix, "generating makefile")

	err := utils.GenerateFile(spec.Paths.Root, "Makefile", makefileTemplate, spec)
	if err != nil {
		return errors.Wrap(err, "failed to create Makefile")
	}

	return nil
}
